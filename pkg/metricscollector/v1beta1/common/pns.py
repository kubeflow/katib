# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import psutil
import sys
import time
import const


def WaitMainProcesses(pool_interval, timout, wait_all, completed_marked_dir):
    """
    Hold metrics collector parser until required pids are finished
    """

    if not sys.platform.startswith('linux'):
        raise Exception("Platform '{}' unsupported".format(sys.platform))

    pids, main_pid = GetMainProcesses(completed_marked_dir)

    return WaitPIDs(pids, main_pid, pool_interval, timout, wait_all, completed_marked_dir)


def GetMainProcesses(completed_marked_dir):
    """
    Return array with all running processes pids and main process pid which metrics collector is waiting.
    """
    pids = set()
    main_pid = 0
    this_pid = psutil.Process().pid

    for proc in psutil.process_iter():
        # Get pid from process object
        pid = proc.pid

        # Get parent process
        ppid = proc.ppid()

        # Ignore the pause container, our own pid, and non-root processes (parent pid != 0)
        if pid == 1 or pid == this_pid or ppid != 0:
            continue

        # Read the process command line, join all cmdlines in one string
        cmd_lind = " ".join(proc.cmdline())

        # By default main_pid is the first process.
        # In addition to that, command line contains completed marker for the main pid.
        # For example: echo completed > /var/log/katib/$$$$.pid
        # completed_marked_dir is the directory for completed marker, e.g. /var/log/katib
        if main_pid == 0 or ("echo {} > {}".format(const.TRAINING_COMPLETED, completed_marked_dir) in cmd_lind):
            main_pid = pid

        pids.add(pid)

    # If mainPid has not been found, return an error.
    if main_pid == 0:
        raise Exception("Unable to find main pid from the process list {}".format(pids))

    return pids, main_pid


def WaitPIDs(pids, main_pid, pool_interval, timout, wait_all, completed_marked_dir):
    """
    Waits until all pids are finished.
    If waitAll == false WaitPIDs waits until main process is finished.
    """
    start = 0
    # not_finished_pids contains pids that are not finished yet
    not_finished_pids = set(pids)

    if pool_interval <= 0:
        raise Exception("Poll interval seconds must be a positive integer")

    # Start main wait loop
    # We should exit when timeout is out or not_finished_pids is empty
    while (timout <= 0 or start < timout) and len(not_finished_pids) > 0:
        finished_pids = set()
        for pid in not_finished_pids:
            # If pid is completed /proc/<pid> dir doesn't exist
            path = "/proc/{}".format(pid)
            if not os.path.exists(path):
                if pid == main_pid:
                    # For main_pid we check if file with "completed" marker exists if completed_marked_dir is set
                    if completed_marked_dir:
                        mark_file = os.path.join(completed_marked_dir, "{}.pid".format(pid))
                        # Check if file contains "completed" marker
                        with open(mark_file) as file_obj:
                            contents = file_obj.read()
                            if contents.strip() != const.TRAINING_COMPLETED:
                                raise Exception(
                                    "Unable to find marker: {} in file: {} with contents: {} for pid: {}".format(
                                        const.TRAINING_COMPLETED, mark_file, contents, pid))
                    # Add main pid to finished pids set
                    finished_pids.add(pid)
                    # Exit loop if wait all is false because main pid is finished
                    if not wait_all:
                        return
                # Add not main pid to finished pids set
                else:
                    finished_pids.add(pid)

        # Update not finished pids set with finished_pids set
        not_finished_pids = not_finished_pids - finished_pids
        # Sleep for pool_interval seconds before next attempt
        time.sleep(pool_interval)
        start = start + pool_interval

    # After main loop not finished pids set should be empty
    if not_finished_pids:
        raise Exception("Timed out waiting for pids to complete")

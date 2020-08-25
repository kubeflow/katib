import os
import psutil
import sys
import time
import const


def WaitMainProcesses(pool_interval, timout, wait_all, completed_marked_dir):
    """
    Hold metrics collector parser until required pids are finished
    """

    if sys.platform != "linux":
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

        # Read the process command line
        cmd_lind = proc.cmdline

        # By default main_pid is the first process.
        # Command line contains completed marker for the main pid
        # For example: echo completed > /var/log/katib/$$$$.pid
        # completed_marked_dir is the directory for completed marker, e.g. /var/log/katib
        if main_pid == 0 or ("echo {} > {}".format(const.TRAINING_COMPLETED, completed_marked_dir) in cmd_lind):
            main_pid = pid

        pids.add(pid)

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
                    # Remove main pid from set with pids
                    not_finished_pids.remove(pid)
                    # Exit loop if wait all is false because main pid is finished
                    if not wait_all:
                        return
                # Remove not main pid from set with pids
                else:
                    not_finished_pids.remove(pid)
        # Sleep for pool_interval seconds before next attempt
        time.sleep(pool_interval)
        start = start + pool_interval

    # After main loop not_finished_pids set should be empty
    if not_finished_pids:
        raise Exception("Timed out waiting for pids to complete")

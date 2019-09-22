import os
import psutil
import time

def GetOtherMainProcesses():
  this_pid = psutil.Process().pid
  pids = set()
  for proc in psutil.process_iter():
    pid = proc.pid
    ppid = proc.ppid()
    if pid == 1 or pid == this_pid or ppid != 0:
      # ignore the pause container, our own pid, and non-root processes
      continue
    pids.add(pid)
  return pids

def WaitPIDs(pids, poll_interval_seconds=1, timeout_seconds=0, is_wait_all=False):
  start = 0
  pids = set(pids)
  if poll_interval_seconds <= 0:
    raise Exception("Poll interval seconds must be a positive integer")
  while (timeout_seconds <= 0 or start < timeout_seconds) and len(pids) > 0:
    stop_pids = set()
    for pid in pids:
      path = "/proc/%d" % pid
      if os.path.isdir(path):
        continue
      else:
        if is_wait_all:
          stop_pids.add(pid)
        else:
          return
    if is_wait_all:
      pids = pids - stop_pids
    time.sleep(poll_interval_seconds)
    start = start + poll_interval_seconds

def WaitOtherMainProcesses(poll_interval_seconds=1, timeout_seconds=0, is_wait_all=False):
  return WaitPIDs(GetOtherMainProcesses(), poll_interval_seconds, timeout_seconds, is_wait_all)

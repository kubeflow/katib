/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	gops "github.com/mitchellh/go-ps"
)

var ErrWaitPidTimeout = fmt.Errorf("Timed out waiting for PID to complete")

type WaitPidsOpts struct {
	PollInterval           time.Duration
	Timeout                time.Duration
	WaitAll                bool
	CompletedMarkedDirPath string
}

func Wait(opts WaitPidsOpts) error {
	pids, err := GetOtherMainProcesses()
	if err != nil {
		return err
	}
	return WaitPIDS(pids, opts)
}

func GetOtherMainProcesses() ([]int, error) {
	pids := []int{}
	allProcs, err := gops.Processes()
	if err != nil {
		return pids, fmt.Errorf("Failed to list processes: %v", err)
	}

	thisPID := os.Getpid()
	for _, proc := range allProcs {
		pid := proc.Pid()
		if pid == 1 || pid == thisPID || proc.PPid() != 0 {
			// ignore the pause container, our own pid, and non-root processes
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func WaitPIDS(pids []int, opts ...WaitPidsOpts) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("Platform '%s' unsupported", runtime.GOOS)
	}
	if len(pids) == 0 {
		return nil
	}
	waitAll := false
	var timeout time.Duration
	var pollInterval = time.Second
	if len(opts) > 0 {
		if opts[0].PollInterval != 0 {
			pollInterval = opts[0].PollInterval
		}
		if opts[0].Timeout != 0 {
			timeout = opts[0].Timeout
		}
		waitAll = opts[0].WaitAll
	}

	finishedPids := []int{}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	var timoutCh <-chan time.Time
	if timeout != 0 {
		timoutCh = time.NewTimer(timeout).C
	}
	for {
		select {
		case <-ticker.C:
			for _, pid := range pids {
				path := fmt.Sprintf("/proc/%d", pid)
				_, err := os.Stat(path)
				if err != nil {
					if os.IsNotExist(err) {
						if opts[0].CompletedMarkedDirPath != "" {
							markFile := filepath.Join(opts[0].CompletedMarkedDirPath, fmt.Sprintf("%d.pid", pid))
							if data, err := ioutil.ReadFile(markFile); err != nil {
								return fmt.Errorf("Process %d hadn't completed: %v", pid, err)
							} else {
								if strings.TrimSpace(string(data)) != TrainingCompleted {
									return fmt.Errorf("Process %d hadn't completed", pid)
								}
							}
						}
						if waitAll {
							finishedPids = append(finishedPids, pid)
							if len(finishedPids) == len(pids) {
								return nil
							}
						} else {
							return nil
						}
					}
					return fmt.Errorf("Fail to check process info: %v", err)
				}
			}
		case <-timoutCh:
			return ErrWaitPidTimeout
		}
	}
}

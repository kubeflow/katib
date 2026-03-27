/*
Copyright 2022 The Kubeflow Authors.

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
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestIsAlreadyCompleted(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string // filename -> content
		expected    bool
		expectError bool
	}{
		{
			name:        "No marker files",
			files:       map[string]string{},
			expected:    false,
			expectError: false,
		},
		{
			name:        "Completed marker exists",
			files:       map[string]string{"123.pid": "completed"},
			expected:    true,
			expectError: false,
		},
		{
			name:        "Early-stopped marker only",
			files:       map[string]string{"123.pid": "early-stopped"},
			expected:    false,
			expectError: false,
		},
		{
			name:        "Empty directory path",
			files:       nil,
			expected:    false,
			expectError: false,
		},
		{
			name: "Multiple pid files with one completed (Istio sidecar scenario)",
			files: map[string]string{
				"45.pid": "completed",
				"52.pid": "early-stopped",
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "Multiple early-stopped files, no completed",
			files: map[string]string{
				"45.pid": "early-stopped",
				"52.pid": "early-stopped",
			},
			expected:    false,
			expectError: false,
		},
		{
			name:        "Completed marker with whitespace",
			files:       map[string]string{"123.pid": "  completed  \n"},
			expected:    true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Create test files
			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			}

			// Test
			var dirPath string
			if tt.files != nil {
				dirPath = tmpDir
			}
			result, err := isAlreadyCompleted(dirPath)

			if (err != nil) != tt.expectError {
				t.Errorf("isAlreadyCompleted(%q) unexpected error: %v", dirPath, err)
			}
			if result != tt.expected {
				t.Errorf("isAlreadyCompleted(%q) = %v, want %v", dirPath, result, tt.expected)
			}
		})
	}
}

func TestWaitMainProcesses(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("WaitMainProcesses only works on linux")
	}

	tests := []struct {
		name        string
		files       map[string]string
		expectError bool
	}{
		{
			name:        "Main process exited but completed marker exists",
			files:       map[string]string{"45.pid": "completed"},
			expectError: false,
		},
		{
			name:        "No process and no marker",
			files:       map[string]string{},
			expectError: true,
		},
		{
			name:        "Main process exited with early-stopped only",
			files:       map[string]string{"45.pid": "early-stopped"},
			expectError: true,
		},
		{
			name: "Sidecar still running but completed marker exists (Istio scenario)",
			files: map[string]string{
				"45.pid": "completed",
				"52.pid": "early-stopped",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			}

			opts := WaitPidsOpts{
				PollInterval:           100 * time.Millisecond,
				Timeout:                1 * time.Second,
				WaitAll:                false,
				CompletedMarkedDirPath: tmpDir,
			}

			err := WaitMainProcesses(opts)

			if tt.expectError && err == nil {
				t.Errorf("WaitMainProcesses() expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("WaitMainProcesses() unexpected error: %v", err)
			}
		})
	}
}

func TestWaitPIDsWithMarkerCheck(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string]string
		expectError bool
	}{
		{
			name:        "Non-existent PID with completed marker exits immediately",
			files:       map[string]string{"99999.pid": "completed"},
			expectError: false,
		},
		{
			name:        "Non-existent PID without marker times out",
			files:       map[string]string{},
			expectError: true,
		},
		{
			name:        "Non-existent PID with early-stopped marker times out",
			files:       map[string]string{"99999.pid": "early-stopped"},
			expectError: true,
		},
		{
			name: "Non-existent PID with sidecar and completed marker exits",
			files: map[string]string{
				"99999.pid": "completed",
				"88888.pid": "early-stopped",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to write test file: %v", err)
				}
			}

			fakePids := map[int]bool{99999: true}
			opts := WaitPidsOpts{
				PollInterval:           100 * time.Millisecond,
				Timeout:                1 * time.Second,
				WaitAll:                false,
				CompletedMarkedDirPath: tmpDir,
			}

			err := WaitPIDs(fakePids, 99999, opts)

			if tt.expectError && err == nil {
				t.Errorf("WaitPIDs() expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("WaitPIDs() unexpected error: %v", err)
			}
		})
	}
}

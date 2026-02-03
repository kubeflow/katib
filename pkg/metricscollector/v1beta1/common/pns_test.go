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
	"testing"
)

func TestIsAlreadyCompleted(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string // filename -> content
		expected bool
	}{
		{
			name:     "No marker files",
			files:    map[string]string{},
			expected: false,
		},
		{
			name:     "Completed marker exists",
			files:    map[string]string{"123.pid": "completed"},
			expected: true,
		},
		{
			name:     "Early-stopped marker only",
			files:    map[string]string{"123.pid": "early-stopped"},
			expected: false,
		},
		{
			name:     "Empty directory path",
			files:    nil,
			expected: false,
		},
		{
			name: "Multiple pid files with one completed (Istio sidecar scenario)",
			files: map[string]string{
				"45.pid": "completed",
				"52.pid": "early-stopped",
			},
			expected: true,
		},
		{
			name: "Multiple early-stopped files, no completed",
			files: map[string]string{
				"45.pid": "early-stopped",
				"52.pid": "early-stopped",
			},
			expected: false,
		},
		{
			name:     "Completed marker with whitespace",
			files:    map[string]string{"123.pid": "  completed  \n"},
			expected: true,
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
			result := isAlreadyCompleted(dirPath)

			if result != tt.expected {
				t.Errorf("isAlreadyCompleted(%q) = %v, want %v", dirPath, result, tt.expected)
			}
		})
	}
}

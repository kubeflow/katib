// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dlkctl v1.0"

// NewCommandVersion return  the version command
func NewCommandVersion() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "display the version of dlkctl",
		Long:  `display the version of dlkctl`,
		Args:  cobra.NoArgs,
		Run:   displayVersion,
	}

	return cmd
}

func displayVersion(cmd *cobra.Command, args []string) {
	fmt.Println(version)
}

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
	"net/http"
	"os"

	"github.com/osrg/dlk/dlkctl/utils"
	"github.com/spf13/cobra"
)

//NewCommandDelLearningTasks generate run cmd
func NewCommandDelLearningTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "learningtask",
		Args:    cobra.ExactArgs(1),
		Short:   "Delete learning task",
		Long:    `Delete learning task`,
		Aliases: []string{"lt"},
		Run:     delLearningTasks,
	}

	//set local flag
	utils.AddNameSpaceFlag(cmd)

	//add subcommand

	return cmd
}

//exec parameter
type delLearningTaskConfig struct {
	params utils.Params
	pf     *PersistentFlags
}

//Main Proceduer of delete learningTask command
func delLearningTasks(cmd *cobra.Command, args []string) {

	//parameter check
	fmt.Println("*** CHECK PARAMS ***")
	dec := delLearningTaskConfig{}
	err := dec.checkParams(cmd, args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("completed")

	fmt.Println("*** exec parameters ***")
	dec.displayParams()

	//send Request DELETE /learningTasks results is stored in array of datastore.LearningTaskInfo
	fmt.Println("*** send Request ***")

	err = dec.sendDelRequest(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("completed")
}

//checkParams check and del exec parameter
func (dec *delLearningTaskConfig) checkParams(cmd *cobra.Command, args []string) error {

	//check and del persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		return err
	}

	//check Flags using common parameter checker
	var params utils.Params
	params, err = utils.CheckFlags(cmd)
	if err != nil {
		return err
	}

	//set config values
	dec.pf = pf
	dec.params = params

	return err
}

//sendDelRequest send del learningTask request and return request results
func (dec *delLearningTaskConfig) sendDelRequest(lt string) error {
	//set url
	url := fmt.Sprintf("http://%s/learningTasks/%s/%s", dec.pf.endpoint, dec.params.Ns, lt)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Printf("failed to create DELETE request: %s\n", err)
		return err
	}

	//send REST API Request
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return err
}

//displayParams display parameters used for sending delete learningTask request
func (dec *delLearningTaskConfig) displayParams() {
	fmt.Printf("| %-30s : %s\n", "Search Namespace", dec.params.Ns)
	fmt.Printf("| %-30s : %s\n", "dlkmanager endpoint", dec.pf.endpoint)
}

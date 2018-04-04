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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/osrg/dlk/dlkmanager/datastore"

	"github.com/osrg/dlk/dlkctl/utils"
	"github.com/spf13/cobra"
)

//NewCommandGetLearningTasks generate run cmd
func NewCommandGetLearningTasks() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "learningtasks",
		Args:    cobra.NoArgs,
		Short:   "Display LearningTasks List",
		Long:    `Display list of learningTasks`,
		Aliases: []string{"lts"},
		Run:     getLearningTasks,
	}

	//set local flag
	utils.AddNameSpaceFlag(cmd)

	//add subcommand

	return cmd
}

//exec parameter
type getLearningTasksConfig struct {
	params utils.Params
	pf     *PersistentFlags
}

//Main Proceduer of get learningTasks command
func getLearningTasks(cmd *cobra.Command, args []string) {

	//parameter check
	fmt.Println("*** CHECK PARAMS ***")
	gec := getLearningTasksConfig{}
	err := gec.checkParams(cmd, args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("completed")
	fmt.Println("*** exec parameters ***")
	gec.displayParams()

	//send Request GET /learningTasks results is stored in array of datastore.LearningTaskInfo
	fmt.Println("*** send Request ***")
	rs, err := gec.sendGetRequest()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//show result
	fmt.Println("*** result ***")
	displayGetLearningTasksResult(rs)

}

//checkParams check and get exec parameter
func (gec *getLearningTasksConfig) checkParams(cmd *cobra.Command, args []string) error {
	var err error

	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err = CheckPersistentFlags()
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
	gec.pf = pf
	gec.params = params

	return err
}

//sendGetRequest send get learningTask request and return request results
func (gec *getLearningTasksConfig) sendGetRequest() ([]datastore.LearningTaskInfo, error) {
	//set url
	url := fmt.Sprintf("http://%s/learningTasks/%s", gec.pf.endpoint, gec.params.Ns)

	//send REST API Request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// get and decode response(json)
	rs := []datastore.LearningTaskInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &rs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return rs, err
}

//displayParams display Result of the command
func displayGetLearningTasksResult(rs []datastore.LearningTaskInfo) {

	//print header
	fmt.Printf("%-40s", "learningTask name")
	fmt.Printf("%-10s", "namespace")
	fmt.Printf("%-5s", "gpu")
	fmt.Printf("%-5s", "NrPS")
	fmt.Printf("%-10s", "NrWorker")
	fmt.Printf("%-25s", "created")
	fmt.Printf("%-10s", "priority")
	fmt.Printf("%-15s", "status")
	fmt.Printf("%-20s\n", "ExecTime")

	//print body
	format := "2006-01-02 15:04 MST"
	for _, i := range rs {
		fmt.Printf("%-40s", i.Name)
		fmt.Printf("%-10s", i.Ns)
		fmt.Printf("%-5d", i.Gpu)
		fmt.Printf("%-5d", i.NrPS)
		fmt.Printf("%-10d", i.NrWorker)
		fmt.Printf("%-25s", i.Created.Format(format))
		fmt.Printf("%-10d", i.Priority)
		fmt.Printf("%-15s", i.State)
		fmt.Printf("%-20s\n", i.ExecTime)
	}
}

//displayParams display parameters used for sending get learningTask request
func (gec *getLearningTasksConfig) displayParams() {
	fmt.Printf("| %-30s : %s\n", "Search Namespace", gec.params.Ns)
	fmt.Printf("| %-30s : %s\n", "dlkmanager endpoint", gec.pf.endpoint)

}

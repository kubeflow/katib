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
	"sort"

	"github.com/osrg/dlk/dlkmanager/datastore"

	"github.com/osrg/dlk/dlkctl/utils"
	"github.com/spf13/cobra"
)

//NewCommandGetLearningTasks generate run cmd
func NewCommandGetLearningTask() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "learningtask",
		Args:    cobra.ExactArgs(1),
		Short:   "Display LearningTask Info",
		Long:    `Display Information of a specified learningTask`,
		Aliases: []string{"lt"},
		Run:     getLearningTask,
	}

	//set local flag
	utils.AddNameSpaceFlag(cmd)

	//add subcommand

	return cmd
}

//exec parameter
type getLearningTaskConfig struct {
	params utils.Params
	pf     *PersistentFlags
}

//Main Proceduer of get learningTasks command
func getLearningTask(cmd *cobra.Command, args []string) {

	//parameter check
	fmt.Println("*** CHECK PARAMS ***")
	gec := getLearningTaskConfig{}
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
	rs, err := gec.sendGetRequest(args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("completed")

	//show result
	fmt.Println("*** result ***")
	if rs.Name == "" {
		fmt.Printf("learningTask %s does not exist\n", args[0])
		return
	}
	displayGetLearningTaskResult(rs)

}

//checkParams check and get exec parameter
func (gec *getLearningTaskConfig) checkParams(cmd *cobra.Command, args []string) error {
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
func (gec *getLearningTaskConfig) sendGetRequest(lt string) (*datastore.LearningTaskInfo, error) {
	//set url
	url := fmt.Sprintf("http://%s/learningTask/%s/%s", gec.pf.endpoint, gec.params.Ns, lt)

	//send REST API Request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// get and decode response(json)
	rs := &datastore.LearningTaskInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, rs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return rs, err
}

//displayParams display Result of the command
func displayGetLearningTaskResult(rs *datastore.LearningTaskInfo) {
	//print body
	format := "2006-01-02 15:04 MST"
	fmt.Printf("| %-30s : %s\n", "learningTask", rs.Name)
	fmt.Printf("| %-30s : %s\n", "NameSpace", rs.Ns)
	fmt.Printf("| %-30s : %d\n", "GPU", rs.Gpu)
	fmt.Printf("| %-30s : %d\n", "NrPS", rs.NrPS)
	fmt.Printf("| %-30s : %d\n", "NrWorker", rs.NrWorker)
	fmt.Printf("| %-30s : %s\n", "Created", rs.Created.Format(format))
	fmt.Printf("| %-30s : %s\n", "Status", rs.State)
	fmt.Printf("| %-30s : %s\n", "PsImage", rs.PsImage)
	fmt.Printf("| %-30s : %s\n", "WorkerImage", rs.WorkerImage)
	fmt.Printf("| %-30s : %s\n", "Scheduler", rs.Scheduler)
	fmt.Printf("| %-30s : %d\n", "Timeout", rs.Timeout)
	fmt.Printf("| %-30s : %s\n", "PVC", rs.Pvc)
	fmt.Printf("| %-30s : %d\n", "Priority", rs.Priority)
	fmt.Printf("| %-30s : %s\n", "ExecTime", rs.ExecTime)

	//pod state
	var keys []string
	for key := range rs.PodState {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	fmt.Printf("| %-30s :\n", "PodState")
	for _, key := range keys {
		fmt.Printf("|  %-40s: %s\n", key, rs.PodState[key])
	}

}

//displayParams display parameters used for sending get learningTask request
func (gec *getLearningTaskConfig) displayParams() {
	fmt.Printf("| %-30s : %s\n", "Search Namespace", gec.params.Ns)
	fmt.Printf("| %-30s : %s\n", "dlkmanager endpoint", gec.pf.endpoint)

}

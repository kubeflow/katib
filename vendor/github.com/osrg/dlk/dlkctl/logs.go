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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/osrg/dlk/dlkctl/utils"
	"github.com/osrg/dlk/dlkmanager/datastore"
	"github.com/spf13/cobra"
)

//LogsConfig :logs exec parameter
type LogsConfig struct {
	params utils.Params
	pf     *PersistentFlags
}

//NewCommandLogs generate logs cmd
func NewCommandLogs() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs learningTaskName [ps|worker|all]",
		Short:   "Display LearningTasks' output Log",
		Long:    `Display LearningTasks' stderr and stdout Log`,
		Aliases: []string{"log"},
		Run:     getLogs,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must specify learningtask and jobType to display logs")
			} else if len(args) == 1 {
				return fmt.Errorf("jobType must be specified for learningTask \"%s\", choose one of: [ps worker all]", args[0])
			} else if len(args) == 2 {
				switch args[1] {
				case "ps":
				case "worker":
				case "all":
					break
				default:
					return fmt.Errorf("argument \"%s\" is not vaild jobType, choose one of: [ps worker all]", args[1])
				}
			} else if len(args) > 2 {
				return errors.New("too may arguments")
			}
			return nil
		},
	}

	//set local flag
	utils.AddNameSpaceFlag(cmd)
	utils.AddSinceTimeFlag(cmd)
	//add subcommand

	return cmd
}

// logs comand entry point
func getLogs(cmd *cobra.Command, args []string) {

	//parameter check
	fmt.Println("*** CHECK PARAMS ***")
	lc := LogsConfig{}
	err := lc.checkParams(cmd, args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("completed")

	fmt.Println("*** GET LOGS ***")

	//sent GET request
	rs, err := lc.sendGetRequest(args[0], args[1])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	//print result
	err = lc.printResult(rs, args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}

//checkParams get and check cli input
func (lc *LogsConfig) checkParams(cmd *cobra.Command, args []string) error {
	var err error

	//check persistent flag
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
	lc.pf = pf
	lc.params = params

	return err
}

//sendGetRequest send get logs request and return request results
func (lc *LogsConfig) sendGetRequest(lt string, jobType string) (*datastore.LtLogInfo, error) {
	//set url
	str := fmt.Sprintf("http://%s/learningTasks/logs/%s/%s/%s", lc.pf.endpoint, lc.params.Ns, lt, jobType)
	reqURL, err := url.Parse(str)
	if err != nil {
		return nil, err
	}

	//if sinceTime is specified,set the value as querry parameter
	if lc.params.SinceTime != "" {
		parameters := url.Values{}
		parameters.Add("sinceTime", lc.params.SinceTime)
		reqURL.RawQuery = parameters.Encode()
	}
	//send REST API Request
	fmt.Println(reqURL.String())
	resp, err := http.Get(reqURL.String())
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// get and decode response(json)
	rs := &datastore.LtLogInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body[:]))
	}
	err = json.Unmarshal(body, rs)
	if err != nil {
		fmt.Println(err.Error())
	}

	return rs, err
}

//printResult print Log request result
func (lc *LogsConfig) printResult(rs *datastore.LtLogInfo, lt string) error {
	//check whether learningTask exist or not
	if rs.LtName == "" {
		return fmt.Errorf("learningTask \"%s\" not found", lt)
	}

	srcLayout := "2006-01-02T15:04:05Z"
	timezone := time.Local

	for i, pod := range rs.PodLogs {
		// Job Name: learningtask-2018-2-6-16-0-52-ps-0
		fmt.Printf("Job Name: %s\n", pod.PodName)

		for _, log := range pod.Logs {
			t, err := time.Parse(srcLayout, log.Time)
			if err != nil {
				fmt.Println(err)
			}

			// [2018-02-06 07:04:14 UTC] starting a server
			fmt.Printf("[%s] %s\n", t.In(timezone).Format(time.RFC3339), log.Value)
		}
		if i != len(rs.PodLogs)-1 {
			fmt.Println("")
		}
	}
	return nil
}

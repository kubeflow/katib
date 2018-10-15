// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"

	"github.com/kubeflow/katib/pkg/api"
)

type pullStudyOpt struct {
	outfile string
	args    []string
}

//NewCommandPullStudy generate pull studies cmd
func NewCommandPullStudy() *cobra.Command {
	var opt pullStudyOpt
	cmd := &cobra.Command{
		Use:   "studies",
		Args:  cobra.ExactArgs(1),
		Short: "Export a Study and its Models lnfo",
		Long:  `Export Information of a Study and its Models to yaml format`,
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			pullStudy(cmd, &opt)
		},
		Aliases: []string{"st"},
	}
	cmd.Flags().StringVarP(&opt.outfile, "output", "o", "", "File path to export")
	return cmd
}

func pullStudy(cmd *cobra.Command, opt *pullStudyOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
	}
	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := api.NewManagerClient(conn)
	listreq := &api.GetStudyListRequest{}
	listr, err := c.GetStudyList(context.Background(), listreq)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
		return
	}
	studyID := ""
	// Search study by Study ID or name
	for _, si := range listr.StudyOverviews {
		if utf8.RuneCountInString(opt.args[0]) >= 7 {
			if strings.HasPrefix(si.Id, opt.args[0]) {
				studyID = si.Id
				break
			}
		}
		if si.Name == opt.args[0] {
			studyID = si.Id
			break
		}
	}
	if studyID == "" {
		log.Fatalf("Study %s is not found", opt.args[0])
	}
	req := &api.GetStudyRequest{
		StudyId: studyID,
	}
	r, err := c.GetStudy(context.Background(), req)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
	}
	mreq := &api.GetSavedModelsRequest{
		StudyName: r.StudyConfig.Name,
	}
	mr, err := c.GetSavedModels(context.Background(), mreq)
	if err != nil {
		log.Fatalf("GetModel failed: %v", err)
	}
	sd := StudyData{
		StudyConf: r.StudyConfig,
		Models:    mr.Models,
	}
	yst, err := yaml.Marshal(sd)
	if err != nil {
		log.Fatalf("Failed to Marshal: %v", err)
	}
	if opt.outfile != "" {
		ioutil.WriteFile(opt.outfile, yst, os.ModePerm)
	} else {
		fmt.Println(string(yst))
	}
}

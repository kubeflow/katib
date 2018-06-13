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

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"

	"github.com/kubeflow/katib/pkg/api"
)

type createStudyOpt struct {
	conf string
	args []string
}

//NewCommandCreateStudy generate create study cmd
func NewCommandCreateStudy() *cobra.Command {
	var opt createStudyOpt
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.NoArgs,
		Short:   "Create a study from a file",
		Long:    "YAML formats are accepted.",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			createStudy(cmd, &opt)
		},
	}
	cmd.Flags().StringVarP(&opt.conf, "config", "f", "", "File path of study config(required)")
	cmd.MarkFlagRequired("config")
	return cmd
}

func createStudy(cmd *cobra.Command, opt *createStudyOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}

	var sc api.StudyConfig
	buf, _ := ioutil.ReadFile(opt.conf)
	err = yaml.Unmarshal(buf, &sc)

	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer conn.Close()
	req := &api.CreateStudyRequest{StudyConfig: &sc}
	c := api.NewManagerClient(conn)
	r, err := c.CreateStudy(context.Background(), req)
	if err != nil {
		log.Fatalf("CreateStudy failed: %v", err)
	}
	fmt.Printf("Study %v is created.", r.StudyId)
}

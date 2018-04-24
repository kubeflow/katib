package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
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

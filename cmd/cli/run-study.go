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

type runStudyOpt struct {
	conf string
	args []string
}

//NewCommandRunStudy generate create study cmd
func NewCommandRunStudy() *cobra.Command {
	var opt runStudyOpt
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.NoArgs,
		Short:   "Run a study from a file",
		Long:    "YAML formats are accepted.",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			runStudy(cmd, &opt)
		},
	}
	cmd.Flags().StringVarP(&opt.conf, "config", "f", "", "File path of study controller config(required)")
	cmd.MarkFlagRequired("config")
	return cmd
}

func runStudy(cmd *cobra.Command, opt *runStudyOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}

	var rsr *api.RunStudyRequest
	buf, err := ioutil.ReadFile(opt.conf)
	if err != nil {
		log.Fatalf("Fail to open %s: %v", opt.conf, err)
		return
	}
	err = yaml.Unmarshal(buf, &rsr)
	if err != nil {
		log.Fatalf("Fail to Unmarshal %s: %v", opt.conf, err)
		return
	}

	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	r, err := c.RunStudy(context.Background(), rsr)
	if err != nil {
		log.Fatalf("RunStudy failed: %v", err)
	}
	fmt.Printf("Study Controller %v is created.\n", r.StudyControllerId)
}

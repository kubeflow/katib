package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

type pushStudyOpt struct {
	file string
	args []string
}

//NewCommandPushStudy generate push model cmd
func NewCommandPushStudy() *cobra.Command {
	var opt pushStudyOpt
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Push a Study Info and its Models from a file or from stdin",
		Long:    "Push a Study Info and its Models from a file or from stdin\nYAML formats are accepted.",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			pushStudy(cmd, &opt)
		},
	}
	cmd.Flags().StringVarP(&opt.file, "file", "f", "", "File path of model config file")
	return cmd
}

func pushStudy(cmd *cobra.Command, opt *pushStudyOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}
	var in StudyData

	if opt.file != "" {
		buf, _ := ioutil.ReadFile(opt.file)
		err = yaml.Unmarshal(buf, &in)
		if err != nil {
			log.Fatalf("Fail to Purse config: %v", err)
			return
		}
	} else if len(opt.args) > 0 {
		err := json.Unmarshal(([]byte)(opt.args[0]), &in)
		if err != nil {
			log.Fatalf("Fail to Purse input: %v", err)
			return
		}
	} else {
		log.Fatalf("You shoud specify study config from a file or options: %v", err)
		return
	}

	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	sreq := &api.CreateStudyRequest{
		StudyConfig: in.StudyConf,
	}
	sr, err := c.CreateStudy(context.Background(), sreq)
	if err != nil {
		log.Fatalf("CreateStudy failed: %v", err)
	}

	for _, m := range in.Models {
		req := &api.SaveModelRequest{
			Model: m,
		}
		_, err = c.SaveModel(context.Background(), req)
		if err != nil {
			log.Fatalf("PushModel failed: %v", err)
		}
		fmt.Printf("Model %v is Pushed.\n", m.WorkerId)
	}
	fmt.Printf("Study %s is Pushd. ID: %s", in.StudyConf.Name, sr.StudyId)
}

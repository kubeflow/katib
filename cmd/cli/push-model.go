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

type pushModelOpt struct {
	file string
	args []string
}

//NewCommandPushModel generate push model cmd
func NewCommandPushModel() *cobra.Command {
	var opt pushModelOpt
	cmd := &cobra.Command{
		Use:     "model",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Push a model Info from a file or from stdin",
		Long:    "YAML or JSON formats are accepted.",
		Aliases: []string{"md"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			pushModel(cmd, &opt)
		},
	}
	cmd.Flags().StringVarP(&opt.file, "file", "f", "", "File path of model config file")
	return cmd
}

func pushModel(cmd *cobra.Command, opt *pushModelOpt) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}
	var req []*api.SaveModelRequest
	if opt.file != "" {
		buf, _ := ioutil.ReadFile(opt.file)
		err = yaml.Unmarshal(buf, &req)
		if err != nil {
			log.Fatalf("Fail to Purse config: %v", err)
			return
		}
	} else if len(opt.args) > 0 {
		err := json.Unmarshal(([]byte)(opt.args[0]), &req)
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
	for _, m := range req {
		_, err = c.SaveModel(context.Background(), m)
		if err != nil {
			log.Fatalf("PushModel failed: %v", err)
		}
		fmt.Printf("Model %v is Pushed.\n", m.Model.WorkerId)
	}
}

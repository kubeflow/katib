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

type pushStudyOpt struct {
	conf  string
	name  string
	owner string
	desc  string
	args  []string
}

//NewCommandGetStudy generate run cmd
func NewCommandPushStudy() *cobra.Command {
	var opt pushStudyOpt
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.NoArgs,
		Short:   "Push a study Info from a file or from option",
		Long:    "YAML or JSON formats are accepted.",
		Aliases: []string{"st"},
		Run: func(cmd *cobra.Command, args []string) {
			opt.args = args
			pushStudy(cmd, &opt)
		},
	}
	cmd.Flags().StringVarP(&opt.conf, "config", "f", "", "File path of study config")
	cmd.Flags().StringVarP(&opt.name, "name", "n", "", "Study name")
	cmd.Flags().StringVarP(&opt.owner, "owner", "o", "Anonymous", "Study owner name")
	cmd.Flags().StringVarP(&opt.desc, "description", "d", "Anonymous", "Description of Study")
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
	var so api.SaveStudyRequest
	if opt.conf != "" {
		var sc api.StudyConfig
		buf, _ := ioutil.ReadFile(opt.conf)
		err = yaml.Unmarshal(buf, &sc)
		if err != nil {
			log.Fatalf("Fail to Purse config: %v", err)
			return
		}
		so.StudyName = sc.Name
		so.Owner = sc.Owner
		so.Description = opt.desc
	} else if opt.name != "" {
		so.StudyName = opt.name
		so.Owner = opt.owner
		so.Description = opt.desc
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
	_, err = c.SaveStudy(context.Background(), &so)
	if err != nil {
		log.Fatalf("PushStudy failed: %v", err)
	}
	fmt.Printf("Study %v is Pushed.\n", so.StudyName)
}

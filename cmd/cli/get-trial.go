package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//NewCommandGetTrial generate run cmd
func NewCommandGetTrial() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "trials",
		Args:    cobra.ExactArgs(1),
		Short:   "Display Trial Info",
		Long:    `Display Information of a trial`,
		Aliases: []string{"tr"},
		Run:     getTrial,
	}
	return cmd
}

func getTrial(cmd *cobra.Command, args []string) {
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()
	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}

	conn, err := grpc.Dial(pf.server, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer conn.Close()

	c := api.NewManagerClient(conn)
	req := &api.GetStudiesRequest{}
	r, err := c.GetStudies(context.Background(), req)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
		return
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintln(w, "StudyID\tName\tOwner\tRunning\tCompleted")
	for _, si := range r.StudyInfos {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n",
			string([]rune(si.StudyId)[:7]),
			si.Name,
			si.Owner,
			si.RunningTrialNum,
			si.CompletedTrialNum)
	}
	fmt.Fprintln(w)
	w.Flush()
}

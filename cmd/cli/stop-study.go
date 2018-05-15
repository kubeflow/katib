package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"unicode/utf8"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//NewCommandStopStudy generate stop study cmd
func NewCommandStopStudy() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.ExactArgs(1),
		Short:   "Stop a study",
		Long:    "Stop study with study ID or study name",
		Aliases: []string{"st"},
		Run:     stopStudy,
	}
	return cmd
}

func stopStudy(cmd *cobra.Command, args []string) {
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
	req := &api.GetStudyListRequest{}
	r, err := c.GetStudyList(context.Background(), req)
	if err != nil {
		log.Fatalf("StopStudy failed: %v", err)
		return
	}
	var sov *api.StudyOverview
	for _, si := range r.StudyOverviews {
		if utf8.RuneCountInString(args[0]) >= 7 {
			if strings.HasPrefix(si.Id, args[0]) {
				sov = si
				break
			}
		}
		if si.Name == args[0] {
			sov = si
			break
		}
	}
	if sov == nil {
		log.Fatalf("Study %s is not found.", args[0])
		return
	}
	sreq := &api.StopStudyRequest{StudyId: sov.Id}
	_, err = c.StopStudy(context.Background(), sreq)
	if err != nil {
		log.Fatalf("StopStudy failed: %v", err)
	}

	fmt.Printf("Study %v has been stopped.\n", sov.Id)
}

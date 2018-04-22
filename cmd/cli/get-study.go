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

//NewCommandGetStudy generate run cmd
func NewCommandGetStudy() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "study",
		Args:    cobra.ExactArgs(1),
		Short:   "Display Study Info",
		Long:    `Display Information of a study`,
		Aliases: []string{"st"},
		Run:     getStudy,
	}
	return cmd
}

func getStudy(cmd *cobra.Command, args []string) {
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
	var sInfo *api.StudyInfo
	for _, si := range r.StudyInfos {
		if utf8.RuneCountInString(args[0]) >= 7 {
			if strings.HasPrefix(si.StudyId, args[0]) {
				sInfo = si
				break
			}
		}
		if si.Name == args[0] {
			sInfo = si
			break
		}
	}
	if sInfo == nil {
		log.Fatalf("Study %s is not found.", args[0])
		return
	}
	fmt.Printf("Study ID:       %s\n", sInfo.StudyId)
	fmt.Printf("Study Name:     %s\n", sInfo.Name)
	fmt.Printf("Study Owner:    %s\n", sInfo.Owner)
	fmt.Printf("Runinng Trial:  %d\n", sInfo.RunningTrialNum)
	fmt.Printf("Complete Trial: %d\n", sInfo.CompletedTrialNum)
}

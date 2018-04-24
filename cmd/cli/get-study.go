package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

//NewCommandGetStudy generate get studies cmd
func NewCommandGetStudy() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "studies",
		Args:    cobra.MaximumNArgs(1),
		Short:   "Display Study lnfo",
		Long:    `Display Information of a studies`,
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
	var sis []*api.StudyInfo
	// Search study if Study ID or name is set
	if len(args) > 0 {
		for _, si := range r.StudyInfos {
			if utf8.RuneCountInString(args[0]) >= 7 {
				if strings.HasPrefix(si.StudyId, args[0]) {
					sis = append(sis, si)
					break
				}
			}
			if si.Name == args[0] {
				sis = append(sis, si)
				break
			}
		}
	} else {
		sis = r.StudyInfos
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintln(w, "StudyID\tName\tOwner\tRunning\tCompleted")
	for _, si := range sis {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%d\n",
			string([]rune(si.StudyId)[:7]),
			si.Name,
			si.Owner,
			si.RunningTrialNum,
			si.CompletedTrialNum)
	}
	w.Flush()
}

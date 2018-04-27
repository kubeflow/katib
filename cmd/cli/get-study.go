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
	req := &api.GetStudyListRequest{}
	r, err := c.GetStudyList(context.Background(), req)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
		return
	}
	result := []*api.StudyOverview{}
	// Search study if Study ID or name is set
	if len(args) > 0 {
		for _, si := range r.StudyOverviews {
			if utf8.RuneCountInString(args[0]) >= 7 {
				if strings.HasPrefix(si.Id, args[0]) {
					result = append(result, si)
					break
				}
			}
			if si.Name == args[0] {
				result = append(result, si)
				break
			}
		}
	} else {
		result = r.StudyOverviews
	}
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintln(w, "StudyID\tName\tOwner")
	for _, si := range result {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			string([]rune(si.Id)[:7]),
			si.Name,
			si.Owner,
		)
	}
	w.Flush()
}

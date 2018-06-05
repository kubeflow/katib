// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/kubeflow/katib/pkg/api"
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

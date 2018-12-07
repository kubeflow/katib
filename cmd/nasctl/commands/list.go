package commands

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kubeflow/katib/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

//NewCommandList generate list cmd
func NewCommandList() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all NAS jobs",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			listNasJobs()
		},
	}

	return cmd
}

func listNasJobs() {

	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()

	if err != nil {
		log.Fatalf("Check persistent flags failed: %v", err)
	}

	//create connection to GRPC server
	connection, err := grpc.Dial(pf.server, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer connection.Close()

	c := api.NewManagerClient(connection)

	getStudyListReq := &api.GetStudyListRequest{}
	getStudyListResp, err := c.GetStudyList(context.Background(), getStudyListReq)

	if err != nil {
		log.Fatalf("GetStudyList failed: %v", err)
	}

	if len(getStudyListResp.StudyOverviews) == 0 {
		log.Fatalf("No Study found")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 14, 0, '\t', tabwriter.TabIndent)
	fmt.Fprintf(w, "NasJobId\tName\n")

	for _, study := range getStudyListResp.StudyOverviews {
		fmt.Fprintf(w, "%v\t%v\n", study.GetId(), study.GetName())
	}
	w.Flush()

}

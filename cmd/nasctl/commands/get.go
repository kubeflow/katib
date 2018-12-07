package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

var (
	outputFormat string
)

//NewCommandGet generate get cmd
func NewCommandGet() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "get NASJobId",
		Short: "Display details about a NAS job",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			getNasJob(args[0], outputFormat)
		},
	}
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: json|yaml")
	cmd.MarkFlagRequired("output")

	return cmd
}

func getNasJob(NASJobID string, outputFormat string) {

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

	getStudyReq := &api.GetStudyRequest{StudyId: NASJobID}

	getStudyResp, err := c.GetStudy(context.Background(), getStudyReq)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
	}

	switch outputFormat {
	case "json":
		outBytes, _ := json.MarshalIndent(getStudyResp, "", "    ")
		fmt.Println(string(outBytes))
	case "yaml":
		outBytes, _ := yaml.Marshal(getStudyResp)
		fmt.Println(string(outBytes))
	default:
		log.Fatalf("Unknown output format: %s", outputFormat)
	}
}

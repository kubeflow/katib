package commands

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	katibClient "github.com/kubeflow/katib/pkg/manager/studyjobclient"
	"github.com/spf13/cobra"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	namespace string
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
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of studyjobs(required)")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func listNasJobs() {

	//Get k8s config
	config := parseKubernetesConfig()

	//Create StudyJobClient
	studyJobClient, err := katibClient.NewStudyjobClient(config)
	if err != nil {
		log.Fatalf("NewStudyJobClient failed: %v", err)
	}

	//Get list of StudyJobs
	studyJobList, err := studyJobClient.GetStudyJobList(namespace)
	if err != nil {
		log.Fatalf("GetStudyJobList failed: %v", err)
	}

	if len(studyJobList.Items) == 0 {
		log.Fatalf("No Study found")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 12, 0, '\t', 0)
	fmt.Fprintln(w, "NasJobId\tName\tStatus\t")

	for _, study := range studyJobList.Items {
		fmt.Fprintf(w, "%v\t%v\t%v\t\n", study.Status.StudyID, study.Spec.StudyName, study.Status.Condition)

	}
	w.Flush()

}

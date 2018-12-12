package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ghodss/yaml"
	apis "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	katibClient "github.com/kubeflow/katib/pkg/manager/studyjobclient"
	"github.com/spf13/cobra"
)

var (
	conf string
)

//NewCommandCreate generate get cmd
func NewCommandCreate() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new NAS job from yaml file",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			createNasJob(conf)
		},
	}

	cmd.Flags().StringVarP(&conf, "config", "f", "", "File path of study config(required)")
	cmd.MarkFlagRequired("config")
	return cmd
}

func createNasJob(conf string) {

	var studyJob apis.StudyJob

	//Parse yaml File
	buf, _ := ioutil.ReadFile(conf)
	err := yaml.Unmarshal(buf, &studyJob)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}

	//Get namespace from studyJob
	namespace := studyJob.ObjectMeta.Namespace

	//Get k8s config
	config := parseKubernetesConfig()

	//Create StudyJobClient
	studyJobClient, err := katibClient.NewStudyjobClient(config)

	if err != nil {
		log.Fatalf("NewStudyJobClient failed: %v", err)
	}

	createStudyJob, err := studyJobClient.CreateStudyJob(&studyJob, namespace)

	if err != nil {
		log.Fatalf("CreateStudyJob failed: %v", err)
	}

	fmt.Printf("NAS job %v is created\n", createStudyJob.ObjectMeta.Name)
}

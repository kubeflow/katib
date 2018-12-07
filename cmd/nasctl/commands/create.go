package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kubeflow/katib/pkg/api"
	apis "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
)

var (
	conf string
)

var OptimizationType_value = map[string]int32{
	"UNKNOWN_OPTIMIZATION": 0,
	"MINIMIZE":             1,
	"MAXIMIZE":             2,
}

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
			createNasJob(conf, args)
		},
	}

	cmd.Flags().StringVarP(&conf, "config", "f", "", "File path of study config(required)")
	cmd.MarkFlagRequired("config")
	return cmd
}

func createNasJob(conf string, args []string) {
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

	var studyJob apis.StudyJob

	buf, _ := ioutil.ReadFile(conf)
	err = yaml.Unmarshal(buf, &studyJob)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}

	var studyJobConfig api.StudyConfig

	studyJobConfig.Name = studyJob.Spec.StudyName
	studyJobConfig.Owner = studyJob.Spec.Owner

	OptimizationTypeString := strings.ToUpper(string(studyJob.Spec.OptimizationType))
	studyJobConfig.OptimizationType = api.OptimizationType(OptimizationType_value[OptimizationTypeString])
	studyJobConfig.OptimizationGoal = *studyJob.Spec.OptimizationGoal
	parameterConfigList := make([]*api.ParameterConfig, 0)
	for _, parameterConfig := range studyJob.Spec.ParameterConfigs {
		append(parameterConfigList, &parameterConfig)
	}

	studyJobConfig.ParameterConfigs = parameterConfigList

	fmt.Println(studyJobConfig.OptimizationType)

	// createStudyReq := &api.CreateStudyRequest{StudyConfig: &studyJobConfig}
	// createStudyResp, err := c.CreateStudy(context.Background(), createStudyReq)
	fmt.Println(c)
	// if err != nil {
	// 	log.Fatalf("CreateStudy failed: %v", err)
	// }

	// fmt.Printf("NAS job %v is created", createStudyResp.StudyId)
}

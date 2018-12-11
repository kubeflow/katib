package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	apis "github.com/kubeflow/katib/pkg/api/operators/apis/studyjob/v1alpha1"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	// yaml "gopkg.in/yaml.v2"
	"github.com/ghodss/yaml"
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

	// c := api.NewManagerClient(connection)

	var studyJob apis.StudyJob

	fmt.Println(conf)
	buf, _ := ioutil.ReadFile(conf)
	// err = json.Unmarshal(buf, &studyJob) works
	err = yaml.Unmarshal(buf, &studyJob)
	if err != nil {
		log.Fatalf("Unmarshal failed: %v", err)
	}
	fmt.Println(studyJob)

	// fmt.Println(studyJob.ObjectMeta)

	// var studyJobConfig api.StudyConfig

	// studyJobConfig.Name = studyJob.Spec.StudyName
	// studyJobConfig.Owner = studyJob.Spec.Owner

	// optimizationTypeString := strings.ToUpper(string(studyJob.Spec.OptimizationType))
	// studyJobConfig.OptimizationType = api.OptimizationType(api.OptimizationType_value[optimizationTypeString])
	// studyJobConfig.OptimizationGoal = *studyJob.Spec.OptimizationGoal

	// var parameterConfigs api.StudyConfig_ParameterConfigs
	// parameterConfigList := make([]*api.ParameterConfig, 0)

	// for _, parameterConfig := range studyJob.Spec.ParameterConfigs {

	// 	parameterTypeString := strings.ToUpper(string(parameterConfig.ParameterType))

	// 	pcFeasible := api.FeasibleSpace{
	// 		List: parameterConfig.Feasible.List,
	// 		Max:  parameterConfig.Feasible.Max,
	// 		Min:  parameterConfig.Feasible.Min,
	// 	}

	// 	pc := api.ParameterConfig{
	// 		Name:          parameterConfig.Name,
	// 		ParameterType: api.ParameterType(api.ParameterType_value[parameterTypeString]),
	// 		Feasible:      &pcFeasible,
	// 	}

	// 	parameterConfigList = append(parameterConfigList, &pc)
	// 	// fmt.Println(parameterConfigList)

	// }
	// parameterConfigs.Configs = parameterConfigList

	// studyJobConfig.ParameterConfigs = &parameterConfigs
	// outBytes, _ := json.MarshalIndent(studyJobConfig, "", "    ")
	// fmt.Println(string(outBytes))
	// fmt.Println(c)

	// createStudyReq := &api.CreateStudyRequest{StudyConfig: &studyJobConfig}
	// createStudyResp, err := c.CreateStudy(context.Background(), createStudyReq)

	// if err != nil {
	// 	log.Fatalf("CreateStudy failed: %v", err)
	// }

	// fmt.Printf("NAS job %v is created", createStudyResp.StudyId)

}

package commands

import (
	"context"
	"encoding/json"
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

	optimizationTypeString := strings.ToUpper(string(studyJob.Spec.OptimizationType))
	studyJobConfig.OptimizationType = api.OptimizationType(api.OptimizationType_value[optimizationTypeString])
	studyJobConfig.OptimizationGoal = *studyJob.Spec.OptimizationGoal

	var parameterConfigs api.StudyConfig_ParameterConfigs
	parameterConfigList := make([]*api.ParameterConfig, 0)

	for _, parameterConfig := range studyJob.Spec.ParameterConfigs {

		parameterTypeString := strings.ToUpper(string(parameterConfig.ParameterType))

		pcFeasible := api.FeasibleSpace{
			List: parameterConfig.Feasible.List,
			Max:  parameterConfig.Feasible.Max,
			Min:  parameterConfig.Feasible.Min,
		}

		pc := api.ParameterConfig{
			Name:          parameterConfig.Name,
			ParameterType: api.ParameterType(api.ParameterType_value[parameterTypeString]),
			Feasible:      &pcFeasible,
		}

		parameterConfigList = append(parameterConfigList, &pc)
		fmt.Println(parameterConfigList)

	}
	parameterConfigs.Configs = parameterConfigList

	studyJobConfig.ParameterConfigs = &parameterConfigs
	outBytes, _ := json.MarshalIndent(studyJobConfig, "", "    ")
	fmt.Println(string(outBytes))
	fmt.Println(c)

	createStudyReq := &api.CreateStudyRequest{StudyConfig: &studyJobConfig}
	createStudyResp, err := c.CreateStudy(context.Background(), createStudyReq)

	if err != nil {
		log.Fatalf("CreateStudy failed: %v", err)
	}

	fmt.Printf("NAS job %v is created", createStudyResp.StudyId)

	// initializeStudy(&studyJob, "kubeflow")
}

// func initializeStudy(instance *apis.StudyJob, ns string) error {
// 	if instance.Spec.SuggestionSpec == nil {
// 		instance.Status.Condition = apis.ConditionFailed
// 		return nil
// 	}
// 	if instance.Spec.SuggestionSpec.SuggestionAlgorithm == "" {
// 		instance.Spec.SuggestionSpec.SuggestionAlgorithm = "random"
// 	}
// 	instance.Status.Condition = apis.ConditionRunning

// 	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
// 	if err != nil {
// 		log.Printf("Connect katib manager error %v", err)
// 		instance.Status.Condition = apis.ConditionFailed
// 		return nil
// 	}
// 	defer conn.Close()
// 	c := api.NewManagerClient(conn)

// 	studyConfig, err := getStudyConf(instance)
// 	if err != nil {
// 		return err
// 	}

// 	log.Printf("Create Study %s", studyConfig.Name)
// 	//CreateStudy
// 	studyID, err := createStudy(c, studyConfig)
// 	if err != nil {
// 		return err
// 	}
// 	instance.Status.StudyID = studyID
// 	log.Printf("Study: %s Suggestion Spec %v", studyID, instance.Spec.SuggestionSpec)
// 	var sspec *apis.SuggestionSpec
// 	if instance.Spec.SuggestionSpec != nil {
// 		sspec = instance.Spec.SuggestionSpec
// 	} else {
// 		sspec = &apis.SuggestionSpec{}
// 	}
// 	sspec.SuggestionParameters = append(sspec.SuggestionParameters,
// 		api.SuggestionParameter{
// 			Name:  "SuggestionCount",
// 			Value: "0",
// 		})
// 	sPID, err := setSuggestionParam(c, studyID, sspec)
// 	if err != nil {
// 		return err
// 	}
// 	instance.Status.SuggestionParameterID = sPID
// 	instance.Status.SuggestionCount += 1
// 	instance.Status.Condition = apis.ConditionRunning
// 	return nil
// }

// func getStudyConf(instance *apis.StudyJob) (*api.StudyConfig, error) {
// 	sconf := &api.StudyConfig{
// 		Metrics: []string{},
// 		ParameterConfigs: &api.StudyConfig_ParameterConfigs{
// 			Configs: []*api.ParameterConfig{},
// 		},
// 	}
// 	sconf.Name = instance.Spec.StudyName
// 	sconf.Owner = instance.Spec.Owner
// 	if instance.Spec.OptimizationGoal != nil {
// 		sconf.OptimizationGoal = *instance.Spec.OptimizationGoal
// 	}
// 	sconf.ObjectiveValueName = instance.Spec.ObjectiveValueName
// 	switch instance.Spec.OptimizationType {
// 	case apis.OptimizationTypeMinimize:
// 		sconf.OptimizationType = api.OptimizationType_MINIMIZE
// 	case apis.OptimizationTypeMaximize:
// 		sconf.OptimizationType = api.OptimizationType_MAXIMIZE
// 	default:
// 		sconf.OptimizationType = api.OptimizationType_UNKNOWN_OPTIMIZATION
// 	}
// 	for _, m := range instance.Spec.MetricsNames {
// 		sconf.Metrics = append(sconf.Metrics, m)
// 	}
// 	for _, pc := range instance.Spec.ParameterConfigs {
// 		p := &api.ParameterConfig{
// 			Feasible: &api.FeasibleSpace{},
// 		}
// 		p.Name = pc.Name
// 		p.Feasible.Min = pc.Feasible.Min
// 		p.Feasible.Max = pc.Feasible.Max
// 		p.Feasible.List = pc.Feasible.List
// 		switch pc.ParameterType {
// 		case apis.ParameterTypeUnknown:
// 			p.ParameterType = api.ParameterType_UNKNOWN_TYPE
// 		case apis.ParameterTypeDouble:
// 			p.ParameterType = api.ParameterType_DOUBLE
// 		case apis.ParameterTypeInt:
// 			p.ParameterType = api.ParameterType_INT
// 		case apis.ParameterTypeDiscrete:
// 			p.ParameterType = api.ParameterType_DISCRETE
// 		case apis.ParameterTypeCategorical:
// 			p.ParameterType = api.ParameterType_CATEGORICAL
// 		}
// 		sconf.ParameterConfigs.Configs = append(sconf.ParameterConfigs.Configs, p)
// 	}
// 	sconf.JobId = string(instance.UID)
// 	return sconf, nil
// }

// func createStudy(c api.ManagerClient, studyConfig *api.StudyConfig) (string, error) {
// 	ctx := context.Background()
// 	createStudyreq := &api.CreateStudyRequest{
// 		StudyConfig: studyConfig,
// 	}
// 	createStudyreply, err := c.CreateStudy(ctx, createStudyreq)
// 	if err != nil {
// 		log.Printf("CreateStudy Error %v", err)
// 		return "", err
// 	}
// 	studyID := createStudyreply.StudyId
// 	log.Printf("Study ID %s", studyID)
// 	getStudyreq := &api.GetStudyRequest{
// 		StudyId: studyID,
// 	}
// 	getStudyReply, err := c.GetStudy(ctx, getStudyreq)
// 	if err != nil {
// 		log.Printf("Study: %s GetConfig Error %v", studyID, err)
// 		return "", err
// 	}
// 	log.Printf("Study ID %s StudyConf%v", studyID, getStudyReply.StudyConfig)
// 	return studyID, nil
// }

// func setSuggestionParam(c api.ManagerClient, studyID string, suggestionSpec *apis.SuggestionSpec) (string, error) {
// 	ctx := context.Background()
// 	pid := ""
// 	if suggestionSpec.SuggestionParameters != nil {
// 		sspr := &api.SetSuggestionParametersRequest{
// 			StudyId:             studyID,
// 			SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
// 		}
// 		for _, p := range suggestionSpec.SuggestionParameters {
// 			sspr.SuggestionParameters = append(
// 				sspr.SuggestionParameters,
// 				&api.SuggestionParameter{
// 					Name:  p.Name,
// 					Value: p.Value,
// 				},
// 			)
// 		}
// 		setSuggesitonParameterReply, err := c.SetSuggestionParameters(ctx, sspr)
// 		if err != nil {
// 			log.Printf("Study %s SetConfig Error %v", studyID, err)
// 			return "", err
// 		}
// 		log.Printf("Study: %s setSuggesitonParameterReply %v", studyID, setSuggesitonParameterReply)
// 		pid = setSuggesitonParameterReply.ParamId
// 	}
// 	return pid, nil
// }

// func getSuggestionParam(c api.ManagerClient, paramID string) ([]*api.SuggestionParameter, error) {
// 	ctx := context.Background()
// 	gsreq := &api.GetSuggestionParametersRequest{
// 		ParamId: paramID,
// 	}
// 	gsrep, err := c.GetSuggestionParameters(ctx, gsreq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return gsrep.SuggestionParameters, err
// }

// func getSuggestion(c api.ManagerClient, studyID string, suggestionSpec *apis.SuggestionSpec, sParamID string) (*api.GetSuggestionsReply, error) {
// 	ctx := context.Background()
// 	getSuggestRequest := &api.GetSuggestionsRequest{
// 		StudyId:             studyID,
// 		SuggestionAlgorithm: suggestionSpec.SuggestionAlgorithm,
// 		RequestNumber:       int32(suggestionSpec.RequestNumber),
// 		//RequestNumber=0 means get all grids.
// 		ParamId: sParamID,
// 	}
// 	getSuggestReply, err := c.GetSuggestions(ctx, getSuggestRequest)
// 	if err != nil {
// 		log.Printf("Study: %s GetSuggestion Error %v", studyID, err)
// 		return nil, err
// 	}
// 	log.Printf("Study: %s CreatedTrials :", studyID)
// 	for _, t := range getSuggestReply.Trials {
// 		log.Printf("\t%v", t)
// 	}
// 	return getSuggestReply, nil
// }

package suggestion

import (
	"context"
	"flag"
	"log"

	"github.com/kubeflow/katib/pkg"
	"github.com/kubeflow/katib/pkg/api"
	katibClient "github.com/kubeflow/katib/pkg/manager/studyjobclient"
	"google.golang.org/grpc"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type CustomSuggestService struct {
}

func NewCustomSuggestService() *CustomSuggestService {
	return &CustomSuggestService{}
}

func (s *CustomSuggestService) GetSuggestions(ctx context.Context, in *api.GetSuggestionsRequest) (*api.GetSuggestionsReply, error) {

	conn, err := grpc.Dial(pkg.ManagerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	defer conn.Close()

	if err != nil {
		log.Fatalf("GetStudyConf failed: %v", err)
		return &api.GetSuggestionsReply{}, err
	}
	reqnum := int(in.RequestNumber)
	sT := make([]*api.Trial, reqnum)

	//Get k8s config
	config := parseKubernetesConfig()
	//Get StudyJobClient
	NewStudyjobClient, err := katibClient.NewStudyjobClient(config)
	if err != nil {
		log.Fatalf("NewStudyjobClient failed: %v", err)
	}

	studyJobClientList, err := NewStudyjobClient.GetStudyJobList()
	log.Println("StudyID=", in.StudyId)
	log.Println("studyJobClientList = ", studyJobClientList)

	return &api.GetSuggestionsReply{Trials: sT}, nil
}

func parseKubernetesConfig() *restclient.Config {
	kubeconfig := flag.String("kubeconfig", "", "Path to a kube config. Only required if out-of-cluster.")
	flag.Parse()

	var err error
	var config *restclient.Config
	if *kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
	} else {
		config, err = restclient.InClusterConfig()
	}
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
}

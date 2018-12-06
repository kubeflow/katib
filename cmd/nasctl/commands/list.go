package commands

import (
	"context"
	"fmt"

	"github.com/kubeflow/katib/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var clientSet = GetKubernetesClient()

//NewCommandList generate list cmd
func NewCommandList() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List NAS jobs",
		Run: func(cmd *cobra.Command, args []string) {

			listNasJobs(args)
		},
	}

	return cmd
}

func listNasJobs(args []string) {

	fmt.Println("This is list command")
	//check and get persistent flag volume
	var pf *PersistentFlags
	pf, err := CheckPersistentFlags()

	if err != nil {
		log.Fatalf("Fail to Check Flags: %v", err)
		return
	}

	connection, err := grpc.Dial(pf.server, grpc.WithInsecure())

	fmt.Println("con = ", connection)

	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer connection.Close()

	c := api.NewManagerClient(connection)
	req := &api.GetStudyListRequest{}

	fmt.Println("c = ", c)
	fmt.Println("req = ", req)
	r, err := c.GetStudyList(context.Background(), req)
	if err != nil {
		log.Fatalf("GetStudy failed: %v", err)
		return
	}

	result := []*api.StudyOverview{}
	fmt.Println("r= ", r.StudyOverviews)
	fmt.Println("result = ", result)

	pods, err := clientSet.CoreV1().Pods("default").List(metav1.ListOptions{})
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
}

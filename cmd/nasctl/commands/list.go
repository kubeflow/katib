package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

var clientSet kubernetes.Interface

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

	if err != nil {
		log.Fatalf("could not connect: %v", err)
		return
	}
	defer connection.Close()

	// pods, err := clientSet.CoreV1().Pods("default").List(metav1.ListOptions{})

	// fmt.Println("pods = ", pods)

}

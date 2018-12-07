package commands

import (
	"errors"
	"flag"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// CLIName is the name of the CLI
	CLIName = "nasctl"
)

//PersistentFlags is used for pass persistent flag value to other command
type PersistentFlags struct {
	server string
}

// NewCommand return a new instance of an nasctl command
func NewCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   CLIName,
		Short: "nasctl is the command line interface for Katib-NAS",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	initFlag(cmd)

	cmd.AddCommand(NewCommandCreate())
	cmd.AddCommand(NewCommandDelete())
	cmd.AddCommand(NewCommandGet())
	cmd.AddCommand(NewCommandList())

	return cmd
}

//CheckPersistentFlags check values is not empty and retun values
func CheckPersistentFlags() (*PersistentFlags, error) {
	var err error

	//katib manager endpoint
	s := viper.GetString("server")
	if s == "" {
		err = errors.New("katib manager endpoint is not specified,use --server provide value")
		return nil, err
	}
	rtn := PersistentFlags{
		server: s,
	}
	return &rtn, err
}

//initFlag manage persistent flags
func initFlag(cmd *cobra.Command) {
	// add Pesistent flags
	cmd.PersistentFlags().StringP("server", "s", "localhost:6789", "katib manager API endpoint")

	//bind viper
	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))
}

func parseKubernetesConfig() *restclient.Config {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	return config
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func GetKubernetesClient() kubernetes.Interface {

	config := parseKubernetesConfig()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("getClusterConfig: %v", err)
	}
	log.Info("Successfully constructed k8s client")
	return client
}

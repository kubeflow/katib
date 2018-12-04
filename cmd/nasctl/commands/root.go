package commands

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cmd.AddCommand(NewCommandDescribe())
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
	//	cobra.OnInitialize(initConfig)
	// add Pesistent flags
	cmd.PersistentFlags().StringP("server", "s", "localhost:6789", "katib manager API endpoint")

	//bind viper
	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))
}

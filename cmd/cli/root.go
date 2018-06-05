package main

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//PersistentFlags is used for pass persistent flag value to other command
type PersistentFlags struct {
	server string
}

// NewRootCommand represents the base command when called without any subcommands
func NewRootCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "katib-cli",
		Short: "katib cli",
		Long:  `This is katib cli client using cobra framework`,
	}

	//initialize config
	initFlag(cmd)

	//add command
	cmd.AddCommand(NewCommandCreate())
	cmd.AddCommand(NewCommandGet())
	cmd.AddCommand(NewCommandStop())
	cmd.AddCommand(NewCommandPush())
	cmd.AddCommand(NewCommandPull())

	//	cmd.AddCommand(NewCommandModel())

	//MISC
	//cmd.AddCommand(NewCommandVersion())

	//Generate bash completion file
	//cmd.AddCommand(NewCommandBashCmp())

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

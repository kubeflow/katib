package main

import (
	"github.com/spf13/cobra"
)

//NewCommandPull generate run cmd
func NewCommandPull() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Pull a resource from a file or from stdin.",
		Long:  `YAML or JSON formats are accepted.`,
	}

	cmd.AddCommand(NewCommandPullStudy())
	//	cmd.AddCommand(NewCommandPullModel())

	return cmd
}

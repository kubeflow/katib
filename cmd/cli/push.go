package main

import (
	"github.com/spf13/cobra"
)

//NewCommandGet generate run cmd
func NewCommandPush() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push a resource from a file or from stdin.",
		Long:  `YAML or JSON formats are accepted.`,
	}

	cmd.AddCommand(NewCommandPushStudy())
	cmd.AddCommand(NewCommandPushModel())

	return cmd
}

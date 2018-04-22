package main

import (
	"github.com/spf13/cobra"
)

//NewCommandGet generate run cmd
func NewCommandCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a resource from a file",
		Long:  `YAML formats are accepted.`,
	}

	cmd.AddCommand(NewCommandCreateStudy())

	return cmd
}

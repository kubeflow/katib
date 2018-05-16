package main

import (
	"github.com/spf13/cobra"
)

//NewCommandRun generate run cmd
func NewCommandRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run study or Trial from a file",
		Long:  `YAML format is accepted.`,
	}

	cmd.AddCommand(NewCommandRunStudy())

	return cmd
}

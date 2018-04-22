package main

import (
	"github.com/spf13/cobra"
)

//NewCommandGet generate run cmd
func NewCommandStop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop a resource",
		Long:  `Specify resource ID or Name.`,
	}

	cmd.AddCommand(NewCommandStopStudy())

	return cmd
}

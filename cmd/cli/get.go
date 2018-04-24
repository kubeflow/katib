package main

import (
	"github.com/spf13/cobra"
)

//NewCommandGet generate get cmd
func NewCommandGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
		Long:  `list of resorces comannd can display includes: studies, study, trials, trial, models, model`,
	}

	//set local flag

	//add subcommand
	cmd.AddCommand(NewCommandGetStudy())
	//	cmd.AddCommand(NewCommandGetTrial())
	cmd.AddCommand(NewCommandGetModel())

	return cmd
}

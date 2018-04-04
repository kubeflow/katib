package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	bcmpfile = "dlkmanager.sh"
)

// NewCommandBashCmp generates bash completion file
func NewCommandBashCmp() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "bash",
		Short: "generate bash completion file",
		Long:  `generate bash completion file, which is copied under /etc/bash_completion.d/`,
		Args:  cobra.NoArgs,
		Run:   genBashCmpFile,
	}

	return cmd
}

func genBashCmpFile(cmd *cobra.Command, args []string) {

	root := cmd.Root()
	root.GenBashCompletionFile(bcmpfile)
	fmt.Println("Bash Completeion File (" + bcmpfile + ") is generated. Please copy it under /etc/bash_completion.d/ and reset your terminal to use autocompletion.")

}

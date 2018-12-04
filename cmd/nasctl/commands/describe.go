package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//NewCommandDescribe generate describe cmd
func NewCommandDescribe() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Describe NAS job",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			desribeNasJob(args[0])
		},
	}

	return cmd
}

func desribeNasJob(args string) {
	fmt.Println("This is describe command")
}

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//NewCommandCreate generate get cmd
func NewCommandCreate() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new NAS job",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			createNasJob(args[0])
		},
	}

	return cmd
}

func createNasJob(args string) {

	fmt.Println("This is create command")
}

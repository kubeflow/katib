package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//NewCommandDelete generate delete cmd
func NewCommandDelete() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete NAS job",
		Run: func(cmd *cobra.Command, args []string) {

			if len(args) != 1 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
			deleteNasJob(args[0])
		},
	}

	return cmd
}

func deleteNasJob(args string) {
	fmt.Println("This is delete command")
}

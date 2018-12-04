package main

import (
	"fmt"
	"os"

	"github.com/kubeflow/katib/cmd/nasctl/commands"
)

//Entry point
func main() {
	//Init command
	if err := commands.NewCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

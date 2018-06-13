package main

import (
	"fmt"
	"os"

	"github.com/kubeflow/katib/cmd/cli/command"
)

//Entry point
func main() {
	//init command
	katibctl := command.NewRootCommand()
	if err := katibctl.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

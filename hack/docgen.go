package main

import (
	"log"

	"github.com/kubeflow/katib/cmd/cli/command"
	"github.com/spf13/cobra/doc"
)

func main() {
	katibCLI := command.NewRootCommand()
	err := doc.GenMarkdownTree(katibCLI, "./docs/CLI")
	if err != nil {
		log.Fatal(err)
	}
}

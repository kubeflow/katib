package main

import (
	"fmt"
	"os"
)

//Entry point
func main() {
	//init command
	katibctl := NewRootCommand()
	if err := katibctl.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

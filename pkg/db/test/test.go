package main

import (
	"fmt"
	"github.com/kubeflow/katib/pkg/db"
	"os"
)

func main() {
	dbInt := db.New()
	study, err := dbInt.GetStudyConfig(os.Args[1])
	if err != nil {
		fmt.Printf("err: %v", err)
	} else {
		fmt.Printf("%v", study)
	}
}

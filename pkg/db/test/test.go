package main

import (
	"fmt"
	"os"

	"github.com/kubeflow/katib/pkg/db"
)

func main() {
	dbInt, err := db.New()
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	study, err := dbInt.GetStudy(os.Args[1])
	if err != nil {
		fmt.Printf("err: %v", err)
	} else {
		fmt.Printf("%v", study)
	}
}

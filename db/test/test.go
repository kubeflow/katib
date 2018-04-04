package main

import (
	"fmt"
	"github.com/mlkube/katib/db"
	"os"
)

func main() {
	db_int := db.New()
	study, err := db_int.GetStudyConfig(os.Args[1])
	if err != nil {
		fmt.Printf("err: %v", err)
	} else {
		fmt.Printf("%v", study)
	}
}

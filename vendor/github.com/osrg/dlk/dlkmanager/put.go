package main

import (
	"fmt"

	"github.com/labstack/echo"
)

func updateLearningTask(c echo.Context) error {
	_, err := fmt.Println("updateLearningTask called")
	return err
}

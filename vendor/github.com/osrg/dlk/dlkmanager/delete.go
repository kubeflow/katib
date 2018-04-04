package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/osrg/dlk/dlkmanager/datastore"
	lgr "github.com/sirupsen/logrus"
)

func deleteLearningTask(c echo.Context) error {

	// get parameter
	ns := c.Param("namespace")
	ltName := c.Param("lt")

	log.WithFields(
		lgr.Fields{
			"API CALLED":  "DELETE /learningTasks",
			"namespace":   ns,
			"lt":          ltName,
			"access from": c.RealIP(),
		}).Info("delete learningTask")

	defer runningLTMu.Unlock()
	runningLTMu.Lock()
	// learning task is running
	if lt, ok := runningLearningTasks[ltName]; ok {

		// notify learning task deletion
		lt.deleteCh <- true
	} else {
		defer completedLTMu.Unlock()
		completedLTMu.Lock()

		// learning task is completed
		if lt, ok = completedLearningTasks[ltName]; ok {
			fmt.Println("learning task is completed")
			return c.NoContent(http.StatusOK)
		} else {
			fmt.Println("learning task is not found within namespace")
			return c.NoContent(http.StatusNotFound)
		}
	}

	// delete LearningTask Info
	err := datastore.Accesor.Remove(ltName)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/osrg/dlk/dlkmanager/datastore"
	lgr "github.com/sirupsen/logrus"
)

func getLearningTasks(c echo.Context) error {

	//get parameter
	ns := c.Param("namespace")

	log.WithFields(
		lgr.Fields{
			"API CALLED":  "GET /learningTasks",
			"namespace":   ns,
			"access from": c.RealIP(),
		}).Info("return learningTasks list")

	//get All LearningTasks
	data, _ := datastore.Accesor.GetAll()

	//filtering by namespace
	rtn := []datastore.LearningTaskInfo{}

	for _, d := range data {
		if d.Ns == ns {
			rtn = append(rtn, d)
		}
	}

	//return result
	return c.JSON(http.StatusOK, rtn)

}

func getLearningTask(c echo.Context) error {

	//get parameter
	ns := c.Param("namespace")
	ltName := c.Param("lt")

	log.WithFields(
		lgr.Fields{
			"API CALLED":  "GET /learningTask",
			"namespace":   ns,
			"lt":          ltName,
			"access from": c.RealIP(),
		}).Info("return learningTask")

	//get LearningTask Info
	rtn, err := datastore.Accesor.Get(ltName)
	if err != nil {
		return c.JSON(http.StatusNotFound, datastore.LearningTaskInfo{})
	}

	//return result
	return c.JSON(http.StatusOK, rtn)
}

func getLearningTaskLogs(c echo.Context) error {

	//get parameter
	ns := c.Param("namespace")
	ltName := c.Param("lt")
	role := c.Param("role")
	sinceTime := c.QueryParam("sinceTime")
	//set default if not specified
	if sinceTime == "" {
		sinceTime = "1700-01-01T01:01:01Z"
	}

	//check time string format RFC3339
	st, err := time.Parse(time.RFC3339, sinceTime)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	log.WithFields(
		lgr.Fields{
			"API CALLED":  "GET /learningTask/logs",
			"namespace":   ns,
			"lt":          ltName,
			"role":        role,
			"sinceTime":   c.QueryParam("sinceTime"),
			"access from": c.RealIP(),
		}).Info("return learningTask logs")

	//filtering by namespace

	var lt *learningTask
	var ok bool

	defer runningLTMu.Unlock()
	runningLTMu.Lock()
	if lt, ok = runningLearningTasks[ltName]; !ok {
		defer completedLTMu.Unlock()
		completedLTMu.Lock()
		if lt, ok = completedLearningTasks[ltName]; !ok {
			return c.JSON(http.StatusNotFound, datastore.LtLogInfo{})
		}
	}

	linfo := datastore.LtLogInfo{LtName: lt.name}
	if role == "worker" || role == "all" {
		for n, logs := range lt.workerLogs {
			lobj := datastore.PodLogInfo{PodName: n}
			for _, l := range logs {
				if l.time.After(st) {
					lobj.Logs = append(lobj.Logs, datastore.LogObj{Time: l.time.Time.Format(time.RFC3339), Value: l.value})
				}
			}
			linfo.PodLogs = append(linfo.PodLogs, lobj)
		}
	}
	if role == "ps" || role == "all" {
		for n, logs := range lt.psLogs {
			lobj := datastore.PodLogInfo{PodName: n}
			for _, l := range logs {
				if l.time.After(st) {
					lobj.Logs = append(lobj.Logs, datastore.LogObj{Time: l.time.Time.Format(time.RFC3339), Value: l.value})
				}
			}
			linfo.PodLogs = append(linfo.PodLogs, lobj)
		}
	}

	//return result
	return c.JSON(http.StatusOK, linfo)

}

package studycontroller

import (
	"context"
	"log"
	"time"

	"github.com/kubeflow/katib/pkg/api"
	"github.com/kubeflow/katib/pkg/db"

	"google.golang.org/grpc"
)

type StudyControllerDefault struct {
	dbIf db.VizierDBInterface
}

func NewStudyControllerDefault() Interface {
	return &StudyControllerDefault{}
}

func (s *StudyControllerDefault) Run(managerAddr string, sctlId string) error {
	s.dbIf = db.New()
	sctl, err := s.dbIf.GetStudyController(sctlId)
	if err != nil {
		return err
	}
	go s.getAndRun(managerAddr, sctl)
	return nil
}

func (s *StudyControllerDefault) getAndRun(managerAddr string, sctl *api.StudyController) {
	conn, err := grpc.Dial(managerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("could not connect: %v", err)
		return
	}
	defer conn.Close()
	c := api.NewManagerClient(conn)
	var updatetimer, estimer *time.Ticker
	updatetimer = time.NewTicker(time.Duration(sctl.UpdateInterval) * time.Second)
	defer updatetimer.Stop()
	if sctl.EarlystoppingInterval > 0 && sctl.EarlystoppingAlgorithm != "" && sctl.EarlystoppingAlgorithm != "none" {
		estimer = time.NewTicker(time.Duration(sctl.EarlystoppingInterval) * time.Second)
		defer estimer.Stop()
	} else {
		//Create dummy Ticker instance
		estimer = time.NewTicker(time.Duration(1) * time.Second)
		estimer.Stop()
	}
	ctx := context.Background()
	suggestReq := &api.GetSuggestionsRequest{
		StudyId:             sctl.StudyId,
		SuggestionAlgorithm: sctl.SuggestionAlgorithm,
		RequestNumber:       sctl.RequestSuggestionNum,
		ParamId:             sctl.SuggestionParamId,
	}
	suggestReply, err := c.GetSuggestions(ctx, suggestReq)
	if err != nil {
		return
	}
	s.dbIf.UpdateStudyControllerState(sctl.StudyControllerId, api.State_RUNNING, "")
	workerIds := []string{}
	for {
		select {
		case <-updatetimer.C:
			var running int32 = 0
			var complete int32 = 0
			getWorkerRequest := &api.GetWorkersRequest{StudyId: sctl.StudyId}
			getWorkerReply, err := c.GetWorkers(ctx, getWorkerRequest)
			if err != nil {
				log.Printf("GetWorker Error %v", err)
				return
			}
			getMetricsRequest := &api.GetMetricsRequest{
				StudyId:   sctl.StudyId,
				WorkerIds: workerIds,
			}
			_, err = c.GetMetrics(ctx, getMetricsRequest)
			if err != nil {
				log.Printf("GetMetrics Error %v", err)
				return
			}
			for _, w := range getWorkerReply.Workers {
				if w.Status == api.State_COMPLETED || w.Status == api.State_KILLED {
					complete++
				} else if w.Status == api.State_RUNNING || w.Status == api.State_PENDING {
					running++
				}
			}
			if complete == sctl.RequestSuggestionNum {
				s.dbIf.UpdateStudyControllerState(sctl.StudyControllerId, api.State_COMPLETED, "")
				return
			}
			if running < sctl.MaxParallel && running+complete < sctl.RequestSuggestionNum {
				reqnum := sctl.MaxParallel - running
				for i := 0; i < int(reqnum); i++ {
					t := suggestReply.Trials[int(complete)+int(running)+i]
					ws := *sctl.WorkerConfig
					for _, p := range t.ParameterSet {
						ws.Command = append(ws.Command, p.Name)
						ws.Command = append(ws.Command, p.Value)
					}
					rtr := &api.RunTrialRequest{
						StudyId:      sctl.StudyId,
						TrialId:      t.TrialId,
						Runtime:      "kubernetes",
						WorkerConfig: &ws,
					}
					runTrialReply, err := c.RunTrial(ctx, rtr)
					if err != nil {
						log.Printf("RunTrial Error %v", err)
						return
					}
					workerIds = append(workerIds, runTrialReply.WorkerId)
				}
			}
		case <-estimer.C:
			//if sctl.EarlystoppingInterval != 0 {
			//}
		}
	}
}

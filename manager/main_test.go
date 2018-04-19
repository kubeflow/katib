package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	api "github.com/kubeflow/hp-tuning/api"
	mockdb "github.com/kubeflow/hp-tuning/mock/db"
	mockworker "github.com/kubeflow/hp-tuning/mock/worker"
)

func TestCreateStudy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mockdb.NewMockVizierDBInterface(ctrl)
	mockWif := mockworker.NewMockWorkerInterface(ctrl)
	sid := "teststudy"
	sc := &api.StudyConfig{
		Name:               "test",
		Owner:              "admin",
		OptimizationType:   1,
		ObjectiveValueName: "obj_name",
		Gpu:                1,
	}
	dbIf = mockDB
	mockDB.EXPECT().CreateStudy(
		sc,
	).Return(sid, nil)

	s := &server{
		wIF:         mockWif,
		StudyChList: make(map[string]studyCh),
	}
	req := &api.CreateStudyRequest{StudyConfig: sc}
	ret, err := s.CreateStudy(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateStudy Error %v", err)
	}
	if ret.StudyId != sid {
		t.Fatalf("Study ID expect "+sid+", get %s", ret.StudyId)
	}
	if len(s.StudyChList) != 1 {
		t.Fatalf("Study register failed. Registered number is %d", len(s.StudyChList))
	} else {
		_, ok := s.StudyChList[sid]
		if !ok {
			t.Fatalf("Study %s is failed to register.", sid)
		}
	}
}
func TestGetStudies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mockdb.NewMockVizierDBInterface(ctrl)
	mockWif := mockworker.NewMockWorkerInterface(ctrl)
	sid := []string{"teststudy1", "teststudy2"}
	s := &server{wIF: mockWif, StudyChList: map[string]studyCh{sid[0]: studyCh{}, sid[1]: studyCh{}}}
	dbIf = mockDB

	sc := []*api.StudyConfig{
		&api.StudyConfig{
			Name:               "test1",
			Owner:              "admin",
			OptimizationType:   1,
			ObjectiveValueName: "obj_name1",
			Gpu:                1,
		},
		&api.StudyConfig{
			Name:               "test2",
			Owner:              "admin",
			OptimizationType:   1,
			ObjectiveValueName: "obj_name2",
		},
	}
	rts := []int32{10, 20}
	cts := []int32{5, 1}
	for i := range sid {
		mockDB.EXPECT().GetStudyConfig(sid[i]).Return(sc[i], nil)
		mockWif.EXPECT().GetRunningTrials(sid[i]).Return(make([]*api.Trial, rts[i]))
		mockWif.EXPECT().GetCompletedTrials(sid[i]).Return(make([]*api.Trial, cts[i]))
	}

	req := &api.GetStudiesRequest{}
	ret, err := s.GetStudies(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateStudy Error %v", err)
	}
	if len(ret.StudyInfos) != len(sid) {
		t.Fatalf("Study Info number %d, expected%d", len(ret.StudyInfos), len(sid))
	} else {
		var j int
		for i := range sid {
			switch ret.StudyInfos[i].StudyId {
			case sid[0]:
				j = 0
			case sid[1]:
				j = 1
			default:
				t.Fatalf("GetStudy Error Study ID %s is not expected", ret.StudyInfos[j].StudyId)
			}
			if ret.StudyInfos[i].Name != sc[j].Name {
				t.Fatalf("GetStudy Error Name %s expected %s", ret.StudyInfos[i].Name, sc[j].Name)
			}
			if ret.StudyInfos[i].Owner != sc[j].Owner {
				t.Fatalf("GetStudy Error Owner %s expected %s", ret.StudyInfos[i].Owner, sc[j].Owner)
			}
			if ret.StudyInfos[i].RunningTrialNum != rts[j] {
				t.Fatalf("GetStudy Error RunningTrialNum %d expected %d", ret.StudyInfos[i].RunningTrialNum, rts[j])
			}
			if ret.StudyInfos[i].CompletedTrialNum != cts[j] {
				t.Fatalf("GetStudy Error CompletedTrialNum %d expected %d", ret.StudyInfos[i].CompletedTrialNum, cts[j])
			}

		}
	}
}

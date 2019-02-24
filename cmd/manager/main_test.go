package main

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	api "github.com/kubeflow/katib/pkg/api"
	mockdb "github.com/kubeflow/katib/pkg/mock/db"
)

func TestCreateStudy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mockdb.NewMockVizierDBInterface(ctrl)
	sid := "teststudy"
	sc := &api.StudyConfig{
		Name:               "test",
		Owner:              "admin",
		OptimizationType:   1,
		ObjectiveValueName: "obj_name",
	}
	dbIf = mockDB
	mockDB.EXPECT().CreateStudy(
		sc,
	).Return(sid, nil)

	s := &server{}

	req := &api.CreateStudyRequest{StudyConfig: sc}
	ret, err := s.CreateStudy(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateStudy Error %v", err)
	}
	if ret.StudyId != sid {
		t.Fatalf("Study ID expect "+sid+", get %s", ret.StudyId)
	}
}
func TestGetStudies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := mockdb.NewMockVizierDBInterface(ctrl)
	sid := []string{"teststudy1", "teststudy2"}
	s := &server{}
	dbIf = mockDB
	sc := []*api.StudyConfig{
		&api.StudyConfig{
			Name:               "test1",
			Owner:              "admin",
			OptimizationType:   1,
			ObjectiveValueName: "obj_name1",
		},
		&api.StudyConfig{
			Name:               "test2",
			Owner:              "admin",
			OptimizationType:   1,
			ObjectiveValueName: "obj_name2",
		},
	}
	mockDB.EXPECT().GetStudyList().Return(sid, nil)
	for i := range sid {
		mockDB.EXPECT().GetStudy(sid[i]).Return(sc[i], nil)
	}

	req := &api.GetStudyListRequest{}
	ret, err := s.GetStudyList(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateStudy Error %v", err)
	}
	if len(ret.StudyOverviews) != len(sid) {
		t.Fatalf("Study Info number %d, expected%d", len(ret.StudyOverviews), len(sid))
	} else {
		var j int
		for i := range sid {
			switch ret.StudyOverviews[i].Id {
			case sid[0]:
				j = 0
			case sid[1]:
				j = 1
			default:
				t.Fatalf("GetStudy Error Study ID %s is not expected", ret.StudyOverviews[j].Id)
			}
			if ret.StudyOverviews[i].Name != sc[j].Name {
				t.Fatalf("GetStudy Error Name %s expected %s", ret.StudyOverviews[i].Name, sc[j].Name)
			}
			if ret.StudyOverviews[i].Owner != sc[j].Owner {
				t.Fatalf("GetStudy Error Owner %s expected %s", ret.StudyOverviews[i].Owner, sc[j].Owner)
			}
		}
	}
}

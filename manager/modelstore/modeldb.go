package modelstore

import (
	"context"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/kubeflow/hp-tuning/api"
	"github.com/kubeflow/hp-tuning/manager/modelstore/modeldb"
	"log"
	"net"
	"strconv"
	"strings"
)

type ModelDB struct {
	host string
	port string
}

func NewModelDB(host string, port string) *ModelDB {
	return &ModelDB{host: host, port: port}
}

func (m *ModelDB) createSocket() (thrift.TTransport, *modeldb.ModelDBServiceClient, error) {
	var trans thrift.TTransport
	var err error
	trans, err = thrift.NewTSocket(net.JoinHostPort(m.host, m.port))
	trans = thrift.NewTFramedTransport(trans)
	if err != nil {
		log.Printf("NewTSocket err %v", err)
		return nil, nil, err
	}
	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	iprot := protocolFactory.GetProtocol(trans)
	oprot := protocolFactory.GetProtocol(trans)
	client := modeldb.NewModelDBServiceClient(thrift.NewTStandardClient(iprot, oprot))
	return trans, client, nil

}

func (m *ModelDB) StoreModel(in *api.StoreModelRequest) error {
	trans, client, err := m.createSocket()
	if err != nil {
		return err
	}
	defer trans.Close()
	if err := trans.Open(); err != nil {
		log.Printf("Error opening socket %v ", err)
		return err
	}
	pids, err := client.GetProjectIds(context.Background(), map[string]string{"Name": in.Model.StudyName})
	if err != nil {
		fmt.Printf("Error get Project IDs %v ", err)
		return err
	}
	if len(pids) == 0 {
		log.Printf("There is no Study name %s, You need to create Study before upload Model.\n", in.Model.StudyName)
	}
	pml, err := m.getProjectModelList(client, in.Model.StudyName)
	if err != nil {
	}
	var did int32 = -1
	var mid int32 = -1
	var msid int32 = -1
	for _, md := range pml {
		if md.Specification.Tag == in.Model.StudyName+":"+in.Model.TrialId {
			log.Printf("Study %s: Trial %s is already exist. Metrics will be updated.\n", in.Model.StudyName, in.Model.TrialId)
			did = md.TrainingDataFrame.ID
			mid = md.ID
			msid = md.Specification.ID
			break
		}
	}

	exrId, err := m.createErun(client, pids[0])
	if err != nil {
		return err
	}
	hs := make([]*modeldb.HyperParameter, len(in.Model.Parameters))
	for i := range in.Model.Parameters {
		hs[i] = &modeldb.HyperParameter{Name: in.Model.Parameters[i].Name, Value: in.Model.Parameters[i].Value}
	}

	dpath := "Unset"
	if in.DataSet != nil {
		dpath = in.DataSet.Path

	}
	df := &modeldb.DataFrame{
		ID:       did,
		Filepath: &dpath,
		NumRows:  -1,
		Metadata: []*modeldb.MetadataKV{},
	}
	did, err = client.StoreDataFrame(context.Background(), df, exrId)
	if err != nil {
		return err
	}
	df.ID = did

	md := &modeldb.Transformer{
		ID:       mid,
		Filepath: &in.Model.ModelPath,
	}
	fe := &modeldb.FitEvent{
		Df:    df,
		Model: md,
		Spec: &modeldb.TransformerSpec{
			ID:              msid,
			TransformerType: "NN",
			Hyperparameters: hs,
			Tag:             in.Model.StudyName + ":" + in.Model.TrialId,
		},
		ExperimentRunId: exrId,
		FeatureColumns:  []string{},
	}
	fres, err := client.StoreFitEvent(context.Background(), fe)
	if err != nil {
		return err
	}
	md.ID = fres.ModelId
	for _, met := range in.Model.Metrics {
		mv, err := strconv.ParseFloat(met.Value, 64)
		if err != nil {
			continue
		}
		me := &modeldb.MetricEvent{
			Df:              df,
			Model:           md,
			MetricType:      met.Name,
			MetricValue:     mv,
			LabelCol:        "label_col",
			PredictionCol:   "prediction_col",
			ExperimentRunId: exrId,
		}
		_, err = client.StoreMetricEvent(context.Background(), me)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *ModelDB) GetStoredStudies() ([]*api.StudyOverView, error) {
	trans, client, err := m.createSocket()
	if err != nil {
		return nil, err
	}
	defer trans.Close()
	if err := trans.Open(); err != nil {
		log.Printf("Error opening socket %v ", err)
		return nil, err
	}
	pov, err := client.GetProjectOverviews(context.Background())
	if err != nil {
		return nil, err
	}

	ret := make([]*api.StudyOverView, len(pov))
	for i := range pov {
		ret[i] = &api.StudyOverView{}
		ret[i].Name = pov[i].Project.Name
		ret[i].Author = pov[i].Project.Author
		ret[i].Description = pov[i].Project.Description
	}
	return ret, nil
}
func (m *ModelDB) convertmdModelToModelInfo(mdm *modeldb.ModelResponse) *api.ModelInfo {
	t := strings.Split(mdm.Specification.Tag, ":")
	var sn, mn string

	if len(t) < 2 {
		sn = "Unknown"
		mn = "Unknown"
	} else {
		sn = t[0]
		mn = t[1]
	}
	param := make([]*api.Parameter, len(mdm.Specification.Hyperparameters))
	for i := range mdm.Specification.Hyperparameters {
		param[i] = &api.Parameter{
			Name:  mdm.Specification.Hyperparameters[i].Name,
			Value: mdm.Specification.Hyperparameters[i].Value,
		}

	}
	met := []*api.Metrics{}
	for k, v := range mdm.Metrics {
		for mk := range v {
			met = append(met, &api.Metrics{
				Name:  k,
				Value: strconv.FormatFloat(v[mk], 'f', 4, 64),
			})
		}
	}
	return &api.ModelInfo{
		StudyName:  sn,
		TrialId:    mn,
		Parameters: param,
		Metrics:    met,
		ModelPath:  *mdm.Metadata,
	}
}

func (m *ModelDB) getProjectModelList(client *modeldb.ModelDBServiceClient, studyName string) ([]*modeldb.ModelResponse, error) {
	pids, err := client.GetProjectIds(context.Background(), map[string]string{"Name": studyName})
	if err != nil {
		log.Printf("Error get Project IDs %v ", err)
		return nil, err
	}
	if len(pids) == 0 {
		log.Printf("Study  %s does not exist.", studyName)
		return nil, err
	}
	pre, err := client.GetRunsAndExperimentsInProject(context.Background(), pids[0])
	if err != nil {
		return nil, err
	}
	var ret []*modeldb.ModelResponse
	for _, er := range pre.ExperimentRuns {
		erd, err := client.GetExperimentRunDetails(context.Background(), er.ID)
		if err != nil {
			return nil, err
		}
		ret = append(ret, erd.ModelResponses...)
	}
	return ret, nil
}

func (m *ModelDB) GetStoredModels(in *api.GetStoredModelsRequest) ([]*api.ModelInfo, error) {
	trans, client, err := m.createSocket()
	if err != nil {
		return nil, err
	}
	defer trans.Close()
	if err := trans.Open(); err != nil {
		log.Printf("Error opening socket %v ", err)
		return nil, err
	}
	pml, err := m.getProjectModelList(client, in.StudyName)
	if err != nil {
		return nil, err
	}
	var ret []*api.ModelInfo
	for _, md := range pml {
		ret = append(ret, m.convertmdModelToModelInfo(md))
	}
	return ret, err
}

func (m *ModelDB) GetStoredModel(in *api.GetStoredModelRequest) (*api.ModelInfo, error) {
	trans, client, err := m.createSocket()
	if err != nil {
		return nil, err
	}
	defer trans.Close()
	if err := trans.Open(); err != nil {
		log.Printf("Error opening socket %v ", err)
		return nil, err
	}
	pml, err := m.getProjectModelList(client, in.StudyName)
	if err != nil {
		return nil, err
	}
	for _, md := range pml {
		if md.Specification.Tag == in.StudyName+":"+in.TrialId {
			return m.convertmdModelToModelInfo(md), nil
		}
	}
	return nil, nil
}

func (m *ModelDB) createErun(client *modeldb.ModelDBServiceClient, pid int32) (int32, error) {
	exe := modeldb.NewExperimentEvent()
	exe.Experiment = modeldb.NewExperiment()
	exe.Experiment.ProjectId = pid
	eres, err := client.StoreExperimentEvent(context.Background(), exe)
	if err != nil {
		log.Printf("StoreExperimentEvent err %v", err)
		return -1, err
	}

	exre := modeldb.NewExperimentRunEvent()
	exre.ExperimentRun = modeldb.NewExperimentRun()
	exre.ExperimentRun.ExperimentId = eres.ExperimentId
	exrres, err := client.StoreExperimentRunEvent(context.Background(), exre)
	if err != nil {
		log.Printf("StoreExperimentRunEvent err %v", err)
		return -1, err
	}
	return exrres.ExperimentRunId, nil
}

func (m *ModelDB) StoreStudy(in *api.StoreStudyRequest) error {
	trans, client, err := m.createSocket()
	if err != nil {
		return err
	}
	defer trans.Close()
	if err := trans.Open(); err != nil {
		log.Printf("Error opening socket %v ", err)
		return err
	}
	pids, err := client.GetProjectIds(context.Background(), map[string]string{"Name": in.StudyName})
	if err != nil {
		log.Printf("Error get Project IDs %v ", err)
		return err
	}

	if len(pids) > 0 {
		log.Printf("Study %s is already exist (Project ID %d)", in.StudyName, pids[0])
	} else {
		pje := modeldb.NewProjectEvent()
		pje.Project = modeldb.NewProject()
		pje.Project.ID = -1
		pje.Project.Name = in.StudyName
		pje.Project.Author = in.Author
		pres, err := client.StoreProjectEvent(context.Background(), pje)
		if err != nil {
			log.Printf("StoreProjectEvent err %v", err)
			return err
		}
		_, err = m.createErun(client, pres.ProjectId)
		if err != nil {
			return err
		}
	}
	return nil
}

package modeldb

import (
	"encoding/json"
	"log"
	"os/exec"
)

type ModelDbReq struct {
	Owner          string             `json:"owner"`
	Study          string             `json:"study"`
	Train          string             `json:"train"`
	ModelPath      string             `json:"modelpath"`
	HyperParameter map[string]string  `json:"hyperparameter"`
	Metrics        map[string]float64 `json:"metrics"`
}

type ModelDbIF struct {
}

func (m *ModelDbIF) SendReq(mr *ModelDbReq) error {
	mr_j, err := json.Marshal(mr)
	if err != nil {
		log.Printf("json marshal err %v", err)
		return err
	}
	out, err := exec.Command("python", "modeldb/Workflow.py", string(mr_j)).CombinedOutput()
	if err != nil {
		log.Printf("exec err %v", string(out))
		return err
	}
	return nil
}

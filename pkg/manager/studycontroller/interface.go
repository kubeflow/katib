package studycontroller

type Interface interface {
	Run(managerAddr string, sctlId string) error
}

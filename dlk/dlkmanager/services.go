package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Generate Service Template
func genSvcTemplate() *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   "", // must be filled by caller
			Labels: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{},
			Ports: []v1.ServicePort{
				v1.ServicePort{
					Protocol:   v1.ProtocolTCP,
					Port:       2222,
					TargetPort: intstr.FromInt(2222),
				},
			},
		},
	}
}

// Generate PS Service
func genPSSvc(name string, ltName string) *v1.Service {
	template := genSvcTemplate()
	template.ObjectMeta.Name = name
	template.ObjectMeta.Labels["learning-task"] = ltName
	template.Spec.Selector["app"] = "tfd-ps"
	return template
}

// Generate Worker Service
func genWorkerSvc(name string, ltName string) *v1.Service {
	template := genSvcTemplate()
	template.ObjectMeta.Name = name
	template.ObjectMeta.Labels["learning-task"] = ltName
	template.Spec.Selector["app"] = "tfd-worker"
	return template
}

// PS Service
type psService struct {
	name string
	svc  *v1.Service
}

func newPSServices(ltName string, nrPS int) []*psService {
	ret := make([]*psService, nrPS)
	for i := 0; i < nrPS; i++ {
		psname := fmt.Sprintf("%s-ps-%d", ltName, i)
		svc := genPSSvc(psname, ltName)
		svc.Spec.Selector["app"] = psname
		ret[i] = &psService{
			name: psname,
			svc:  svc,
		}
	}
	return ret
}

// Worker Service
type workerService struct {
	name string
	svc  *v1.Service
}

func newWorkerServices(ltName string, nrWorker int) []*workerService {
	ret := make([]*workerService, nrWorker)
	for i := 0; i < nrWorker; i++ {
		workername := fmt.Sprintf("%s-worker-%d", ltName, i)
		svc := genWorkerSvc(workername, ltName)
		svc.Spec.Selector["app"] = workername
		ret[i] = &workerService{
			name: workername,
			svc:  svc,
		}
	}
	return ret
}

package tensorboard

import (
	"bytes"
	"io/ioutil"

	apiv1 "k8s.io/api/core/v1"
	exbeatav1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

func initk8sCl() (*kubernetes.Clientset, error) {
	// Create the kubernetes client
	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func SpawnTensorBoard(tid string, studyname string, namespace string, inhost *string, pvc string, modelpath string) error {
	BUFSIZE := 1024
	var tFile []byte
	var err error

	dep := exbeatav1.Deployment{}
	tFile, err = ioutil.ReadFile("/tensorboard/manifest_template/deployment.yaml")
	if err != nil {
		return err
	}
	k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(tFile), BUFSIZE).Decode(&dep)

	ing := exbeatav1.Ingress{}
	tFile, err = ioutil.ReadFile("/tensorboard/manifest_template/ingress.yaml")
	if err != nil {
		return err
	}
	k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(tFile), BUFSIZE).Decode(&ing)

	svc := apiv1.Service{}
	tFile, err = ioutil.ReadFile("/tensorboard/manifest_template/service.yaml")
	if err != nil {
		return err
	}
	k8syaml.NewYAMLOrJSONDecoder(bytes.NewReader(tFile), BUFSIZE).Decode(&svc)

	tname := "tensorboard-" + tid

	dep.ObjectMeta.Name = tname
	dep.ObjectMeta.Labels["TrialID"] = tid
	dep.Spec.Template.ObjectMeta.Labels["TrialID"] = tid
	dep.Spec.Template.Spec.Containers[0].Args = append(dep.Spec.Template.Spec.Containers[0].Args, "--logdir="+modelpath)
	dep.Spec.Template.Spec.Containers[0].Args = append(dep.Spec.Template.Spec.Containers[0].Args, "--path_prefix=/tensorboard/"+studyname+"/"+tid)
	dep.Spec.Template.Spec.Volumes = append(dep.Spec.Template.Spec.Volumes, apiv1.Volume{
		Name: "pvc-mount-point",
		VolumeSource: apiv1.VolumeSource{
			PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
				ClaimName: pvc,
			},
		},
	},
	)
	dep.Spec.Template.Spec.Containers[0].VolumeMounts = append(dep.Spec.Template.Spec.Containers[0].VolumeMounts, apiv1.VolumeMount{
		Name:      "pvc-mount-point",
		MountPath: pvc,
	},
	)

	svc.ObjectMeta.Name = tname
	svc.ObjectMeta.Labels["TrialID"] = tid
	svc.Spec.Selector["TrialID"] = tid

	ing.ObjectMeta.Name = tname
	ing.Spec.Rules[0].Host = *inhost
	ing.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Path = "/tensorboard/" + studyname + "/" + tid
	ing.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.ServiceName = tname

	kcl, _ := initk8sCl()
	_, err = kcl.ExtensionsV1beta1().Deployments(namespace).Create(&dep)
	if err != nil {
		return err
	}
	_, err = kcl.CoreV1().Services(namespace).Create(&svc)
	if err != nil {
		return err
	}
	_, err = kcl.ExtensionsV1beta1().Ingresses(namespace).Create(&ing)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTensorBoard(tid string, namespace string) error {
	tname := "tensorboard-" + tid
	kcl, _ := initk8sCl()
	var err error
	err = kcl.ExtensionsV1beta1().Deployments(namespace).Delete(tname, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = kcl.CoreV1().Services(namespace).Delete(tname, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	err = kcl.ExtensionsV1beta1().Ingresses(namespace).Delete(tname, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

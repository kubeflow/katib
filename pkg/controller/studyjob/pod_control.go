/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package studyjob

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type PodControl struct {
	kubeClient clientset.Interface
}

func NewPodControl() (*PodControl, error) {
	config, err := restclient.InClusterConfig()
	if err != nil {
		return nil, err
	}
	kc, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &PodControl{
		kubeClient: kc,
	}, nil
}

func (c PodControl) DeletePodsForWorker(namespace string, wid string) error {
	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{"job-name":wid},
	})
	if err != nil {
		return err
	}
	log.Printf("Deleting pods with selector %v", selector.String())
	listOptions := metav1.ListOptions{
		LabelSelector: selector.String(),
	}
	err = c.kubeClient.CoreV1().Pods(namespace).DeleteCollection(&metav1.DeleteOptions{}, listOptions)
	if err != nil {
		return err
	}
	log.Printf("Deleted pods with selector %v.", selector.String())
	return nil
}

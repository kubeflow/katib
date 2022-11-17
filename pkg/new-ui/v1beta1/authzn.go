package v1beta1

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
	"github.com/pkg/errors"
	v1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ENV variables
var (
	USER_HEADER  = env.GetEnvOrDefault("USERID_HEADER", "kubeflow-userid")
	USER_PREFIX  = env.GetEnvOrDefault("USERID_PREFIX", ":")
	DISABLE_AUTH = env.GetEnvOrDefault("APP_DISABLE_AUTH", "false") == "true"
	BACKEND_MODE = env.GetEnvOrDefault("BACKEND_MODE", "prod")
)

func GetUsername(r *http.Request) (string, error) {
	var username string
	if DISABLE_AUTH {
		log.Printf("APP_DISABLE_AUTH set to True. Skipping authorization check")
		return "", nil
	}

	if r.Header.Get(USER_HEADER) == "" {
		return "", errors.New("User header not present!")
	}

	user := r.Header.Get(USER_HEADER)
	username = strings.Replace(user, USER_PREFIX, "", 1)

	return username, nil
}

// Function for constructing SubjectAccessReviews (SAR) objects
func CreateSAR(user, verb, namespace, group,
	version, resource, subresource, name string) *v1.SubjectAccessReview {

	sar := &v1.SubjectAccessReview{
		Spec: v1.SubjectAccessReviewSpec{
			User: user,
			ResourceAttributes: &v1.ResourceAttributes{
				Namespace:   namespace,
				Verb:        verb,
				Group:       group,
				Version:     version,
				Resource:    resource,
				Subresource: subresource,
				Name:        name,
			},
		},
	}
	return sar
}

func IsAuthorized(user, verb, namespace, group,
	version, resource, subresource, name string, client *kubernetes.Clientset) error {

	// Skip authz when in dev_mode
	if BACKEND_MODE == "dev" || BACKEND_MODE == "development" {
		log.Printf("Skipping authorization check in development mode")
		return nil
	}
	// Skip authz when admin is explicity requested it
	if DISABLE_AUTH {
		log.Printf("APP_DISABLE_AUTH set to True. Skipping authorization check")
		return nil
	}

	sar := CreateSAR(user, verb, namespace, group, version, resource, subresource, name)

	res, err := client.AuthorizationV1().SubjectAccessReviews().Create(context.TODO(), sar, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error submitting SubjectAccessReview: %v, %s", sar, err.Error())
		return err
	}

	if res.Status.Allowed {
		return nil
	}

	msg := generateUnauthorizedMessage(user, verb, namespace, group, version, resource, subresource, res)
	err = errors.New(msg)
	return err
}

func generateUnauthorizedMessage(user, verb, namespace, group,
	version, resource, subresource string, sar *v1.SubjectAccessReview) string {

	msg := fmt.Sprintf("User: %s is not authorized to %s", user, verb)

	if group == "" {
		msg += fmt.Sprintf(" %s/%s", version, resource)
	} else {
		msg += fmt.Sprintf(" %s/%s/%s", group, version, resource)
	}

	if subresource != "" {
		msg += fmt.Sprintf("/%s", subresource)
	}

	if namespace != "" {
		msg += fmt.Sprintf(" in namespace: %s", namespace)
	}
	if sar.Status.Reason != "" {
		msg += fmt.Sprintf(" ,reason: %s", sar.Status.Reason)
	}
	return msg
}

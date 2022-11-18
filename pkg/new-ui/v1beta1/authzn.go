package v1beta1

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/kubeflow/katib/pkg/util/v1beta1/env"
	v1 "k8s.io/api/authorization/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ENV variables
var (
	USER_HEADER  = env.GetEnvOrDefault("USERID_HEADER", "kubeflow-userid")
	USER_PREFIX  = env.GetEnvOrDefault("USERID_PREFIX", ":")
	DISABLE_AUTH = env.GetEnvOrDefault("APP_DISABLE_AUTH", "false")
)

func GetUsername(r *http.Request) (string, error) {
	var username string
	if DISABLE_AUTH == "true" {
		log.Printf("APP_DISABLE_AUTH set to True. Skipping authentication check")
		return "", nil
	}

	if r.Header.Get(USER_HEADER) == "" {
		return "", errors.New("user header not present")
	}

	user := r.Header.Get(USER_HEADER)
	username = strings.Replace(user, USER_PREFIX, "", 1)

	return username, nil
}

// Function for constructing SubjectAccessReviews (SAR) objects
func CreateSAR(user, verb, namespace, resource, subresource, name string, schema schema.GroupVersion) *v1.SubjectAccessReview {

	sar := &v1.SubjectAccessReview{
		Spec: v1.SubjectAccessReviewSpec{
			User: user,
			ResourceAttributes: &v1.ResourceAttributes{
				Namespace:   namespace,
				Verb:        verb,
				Group:       schema.Group,
				Version:     schema.Version,
				Resource:    resource,
				Subresource: subresource,
				Name:        name,
			},
		},
	}
	return sar
}

func IsAuthorized(user, verb, namespace, resource, subresource, name string, schema schema.GroupVersion, client client.Client) error {

	// Skip authz when admin is explicity requested it
	if DISABLE_AUTH == "true" {
		log.Printf("APP_DISABLE_AUTH set to True. Skipping authorization check")
		return nil
	}

	sar := CreateSAR(user, verb, namespace, resource, subresource, name, schema)

	err := client.Create(context.TODO(), sar)
	if err != nil {
		log.Printf("Error submitting SubjectAccessReview: %v, %s", sar, err.Error())
		return err
	}

	if sar.Status.Allowed {
		return nil
	}

	msg := generateUnauthorizedMessage(user, verb, namespace, resource, subresource, schema, sar)
	return errors.New(msg)
}

func generateUnauthorizedMessage(user, verb, namespace, resource, subresource string, schema schema.GroupVersion, sar *v1.SubjectAccessReview) string {

	msg := fmt.Sprintf("User: %s is not authorized to %s", user, verb)

	if schema.Group == "" {
		msg += fmt.Sprintf(" %s/%s", schema.Version, resource)
	} else {
		msg += fmt.Sprintf(" %s/%s/%s", schema.Group, schema.Version, resource)
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

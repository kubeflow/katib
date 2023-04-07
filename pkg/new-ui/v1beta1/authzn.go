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
	USER_HEADER   = env.GetEnvOrDefault("USERID_HEADER", "kubeflow-userid")
	GROUPS_HEADER = env.GetEnvOrDefault("GROUPS_HEADER", "kubeflow-groups")
	USER_PREFIX   = env.GetEnvOrDefault("USERID_PREFIX", ":")
	DISABLE_AUTH  = env.GetEnvOrDefault("APP_DISABLE_AUTH", "true")
)

// Function for constructing SubjectAccessReviews (SAR) objects
func CreateSAR(user string, userGroups []string, verb, namespace, resource, subresource, name string, schema schema.GroupVersion) *v1.SubjectAccessReview {

	sar := &v1.SubjectAccessReview{
		Spec: v1.SubjectAccessReviewSpec{
			User:   user,
			Groups: userGroups,
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

func IsAuthorized(verb, namespace, resource, subresource, name string, schema schema.GroupVersion, client client.Client, r *http.Request) (string, error) {

	// We disable authn/authz checks when in standalone mode.
	if DISABLE_AUTH == "true" {
		log.Printf("APP_DISABLE_AUTH set to True. Skipping authentication/authorization checks")
		return "", nil
	}
	// Check if an incoming request is from an authenticated user (kubeflow mode: kubeflow-userid header)
	if r.Header.Get(USER_HEADER) == "" {
		return "", errors.New("user header not present")
	}

	user := r.Header.Get(USER_HEADER)
	user = strings.Replace(user, USER_PREFIX, "", 1)

	var groups []string
	if header := r.Header.Get(GROUPS_HEADER); header != "" {
		groups = strings.Split(header, ",")
	}

	// Check if the user is authorized to perform a given action on katib/k8s resources.
	sar := CreateSAR(user, groups, verb, namespace, resource, subresource, name, schema)
	err := client.Create(context.TODO(), sar)
	if err != nil {
		log.Printf("Error submitting SubjectAccessReview: %v, %s", sar, err.Error())
		return user, err
	}

	if sar.Status.Allowed {
		return user, nil
	}

	msg := generateUnauthorizedMessage(user, groups, verb, namespace, resource, subresource, schema, sar)
	return user, errors.New(msg)
}

func generateUnauthorizedMessage(user string, userGroups []string, verb, namespace, resource, subresource string, schema schema.GroupVersion, sar *v1.SubjectAccessReview) string {

	msg := fmt.Sprintf("User: %s with groups %v is not authorized to %s", user, userGroups, verb)

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

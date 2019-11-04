package katibconfig

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	corev1 "k8s.io/api/core/v1"
	apitypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	common "github.com/kubeflow/katib/pkg/apis/controller/common/v1alpha3"
	"github.com/kubeflow/katib/pkg/controller.v1alpha3/consts"
)

func GetSuggestionContainerImage(algorithmName string, client client.Client) (string, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return "", err
	}
	if config, ok := configMap.Data[consts.LabelSuggestionTag]; ok {
		suggestionConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(config), &suggestionConfig); err != nil {
			return "", err
		}
		if imageConfig, ok := suggestionConfig[algorithmName]; ok {
			if image, yes := imageConfig[consts.LabelSuggestionImageTag]; yes {
				if strings.TrimSpace(image) != "" {
					return image, nil
				} else {
					return "", errors.New("Required value for " + consts.LabelSuggestionImageTag + " configuration of algorithm name " + algorithmName)
				}
			} else {
				return "", errors.New("Failed to find " + consts.LabelSuggestionImageTag + " configuration of algorithm name " + algorithmName)
			}
		} else {
			return "", errors.New("Failed to find algorithm image mapping " + algorithmName)
		}
	} else {
		return "", errors.New("Failed to find algorithm image mapping in configmap " + consts.KatibConfigMapName)
	}
}

func GetMetricsCollectorImage(cKind common.CollectorKind, client client.Client) (string, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(
		context.TODO(),
		apitypes.NamespacedName{Name: consts.KatibConfigMapName, Namespace: consts.DefaultKatibNamespace},
		configMap)
	if err != nil {
		return "", err
	}
	if mcs, ok := configMap.Data[consts.LabelMetricsCollectorSidecar]; ok {
		kind := string(cKind)
		mcsConfig := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(mcs), &mcsConfig); err != nil {
			return "", err
		}
		if mc, ok := mcsConfig[kind]; ok {
			if image, yes := mc[consts.LabelMetricsCollectorSidecarImage]; yes {
				if strings.TrimSpace(image) != "" {
					return image, nil
				} else {
					return "", errors.New("Required value for " + consts.LabelMetricsCollectorSidecarImage + "configuration of metricsCollector kind " + kind)
				}
			} else {
				return "", errors.New("Failed to find " + consts.LabelMetricsCollectorSidecarImage + " configuration of metricsCollector kind " + kind)
			}
		} else {
			return "", errors.New("Cannot support metricsCollector injection for kind " + kind)
		}
	} else {
		return "", errors.New("Failed to find metrics collector configuration in configmap " + consts.KatibConfigMapName)
	}
}

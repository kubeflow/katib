`Experiment` is user level CRD. When you want to tune hyperparameter for your machine learning model by Katib, you just need create an Experiment CR. There are several examples [here](../examples/v1alpha3/) for your reference. In this document, fields of `Experiment` will be introduced in detail.

## metadata
- name
- namespace
## spec
- (TODO)
- [spec.metricsCollectorSpec.source.filter.metricsFormat](https://github.com/kubeflow/katib/blob/master/pkg/apis/controller/common/v1alpha3/common_types.go#L94-L98): Its type is string array and each element is regular expression (here it is required that the regular expression must have two top subexpressions, the first matched one will be taken as metric name and the second one as metric value). For now, when metrics collector kind is `File` and `StdOut`, this field can take effect and a default one `([\w|-]+)\s*=\s*((-?\d+)(\.\d+)?)` will be applied if not specified. [Here](https://github.com/kubeflow/katib/blob/master/examples/v1alpha3/file-metricscollector-example.yaml) is an example of customized metrics filter.

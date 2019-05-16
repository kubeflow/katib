FROM tensorflow/tensorflow:1.11.0
RUN pip install rfc3339 grpcio googleapis-common-protos
ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/tfevent-metricscollector/v1alpha1
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/api/v1alpha1/python:/usr/src/app/github.com/kubeflow/katib/pkg/manager/v1alpha1/file-metricscollector/tf-event

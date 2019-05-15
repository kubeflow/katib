FROM tensorflow/tensorflow:1.11.0
RUN pip install rfc3339 grpcio googleapis-common-protos
ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/tfevent-metricscollector/v1alpha2
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/api/v1alpha2/python:/usr/src/app/github.com/kubeflow/katib/pkg/util/v1alpha2/tfevent-metricscollector

FROM tensorflow/tensorflow:1.11.0
RUN pip install rfc3339 grpcio googleapis-common-protos
ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/metricscollector/v1alpha3/tfevent-metricscollector/
RUN pip install --no-cache-dir -r requirements.txt
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/v1alpha3/python:/usr/src/app/github.com/kubeflow/katib/pkg/metricscollector/v1alpha3/tfevent-metricscollector/:/usr/src/app/github.com/kubeflow/katib/pkg/metricscollector/v1alpha3/common/
ENTRYPOINT ["python", "main.py"]

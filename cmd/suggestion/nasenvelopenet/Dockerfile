FROM python:3.6

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/nasenvelopenet/v1alpha1
RUN pip install --no-cache-dir -r requirements.txt
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/api/v1alpha1/python

ENTRYPOINT ["python", "-u", "main.py"]

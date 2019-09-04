FROM python:3

ADD . /usr/src/app/github.com/kubeflow/katib
WORKDIR /usr/src/app/github.com/kubeflow/katib/cmd/suggestion/grid/v1alpha3
RUN if [ "$(uname -m)" = "ppc64le" ]; then \
        apt-get -y update && \
        apt-get -y install gfortran libopenblas-dev liblapack-dev && \
        pip install cython; \
    fi
RUN pip install --no-cache-dir -r requirements.txt
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/katib:/usr/src/app/github.com/kubeflow/katib/pkg/apis/manager/v1alpha3/python

ENTRYPOINT ["python", "main.py"]

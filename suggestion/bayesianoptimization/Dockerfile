FROM python:3

ADD . /usr/src/app/github.com/kubeflow/hp-tuning
WORKDIR /usr/src/app/github.com/kubeflow/hp-tuning/suggestion/bayesianoptimization
RUN pip install --no-cache-dir -r requirements.txt
ENV PYTHONPATH /usr/src/app/github.com/kubeflow/hp-tuning

ENTRYPOINT ["python", "main.py"]
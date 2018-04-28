FROM gcr.io/kubeflow-images-staging/tensorflow-1.6.0-notebook-cpu
COPY requirements.txt /workdir/requirements.txt
RUN pip install -r /workdir/requirements.txt
RUN pip install google.cloud sklearn ktext annoy nltk
COPY train.py /workdir/train.py
COPY seq2seq_utils.py /workdir/seq2seq_utils.py
COPY IssueSummarization.py /workdir/IssueSummarization.py

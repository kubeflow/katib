FROM pytorch/pytorch:0.4_cuda9_cudnn7

RUN mkdir -p /opt/mlkube
COPY dist_sendrecv.py /opt/mlkube/
ENTRYPOINT ["python", "/opt/mlkube/dist_sendrecv.py"]

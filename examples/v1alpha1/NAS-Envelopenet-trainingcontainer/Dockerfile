ARG cuda_version=9.0
ARG cudnn_version=7
FROM nvidia/cuda:${cuda_version}-cudnn${cudnn_version}-devel

# Install system packages
RUN apt-get update && apt-get install -y software-properties-common && \
      add-apt-repository ppa:deadsnakes/ppa && \
      apt-get update && \
      apt-get install -y --no-install-recommends \
      bzip2 \
      g++ \
      git \
      graphviz \
      libgl1-mesa-glx \
      libhdf5-dev \
      openmpi-bin \
      python3.5 \
      python3-pip \
      python3-setuptools \
      python3-dev \
      wget && \
    rm -rf /var/lib/apt/lists/*


ADD . /app
WORKDIR /app

RUN pip3 install --upgrade pip
RUN pip3 install --no-cache-dir -r requirements.txt
ENV PYTHONPATH /app

ENTRYPOINT ["python3.5", "-u", "run_trial.py"]

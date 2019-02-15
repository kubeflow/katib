FROM nvidia/cuda:9.0-cudnn7-devel-ubuntu16.04
ARG PYTHON_VERSION=3.6

RUN apt-get update && apt-get install -y --no-install-recommends \
         build-essential \
         cmake \
         git \
         curl \
         vim \
         ca-certificates \
         openssh-client \
         libjpeg-dev \
         libpng-dev &&\
     rm -rf /var/lib/apt/lists/*

RUN curl -o ~/miniconda.sh -O  https://repo.continuum.io/miniconda/Miniconda3-latest-Linux-x86_64.sh  && \
     chmod +x ~/miniconda.sh && \
     ~/miniconda.sh -b -p /opt/conda && \
     rm ~/miniconda.sh && \
     /opt/conda/bin/conda install -y python=$PYTHON_VERSION numpy pyyaml scipy ipython mkl mkl-include cython typing && \
     /opt/conda/bin/conda install -y -c pytorch magma-cuda90 && \
     /opt/conda/bin/conda install -c conda-forge openmpi && \
     /opt/conda/bin/conda clean -ya
ENV PATH /opt/conda/bin:$PATH
# This must be done before pip so that requirements.txt is available
WORKDIR /opt/pytorch

RUN git clone --recursive https://github.com/pytorch/pytorch

RUN TORCH_CUDA_ARCH_LIST="3.5 5.2 6.0 6.1 7.0+PTX" TORCH_NVCC_FLAGS="-Xfatbin -compress-all" \
    CMAKE_PREFIX_PATH="$(dirname $(which conda))/../" \
    cd pytorch/ && \
    pip install -v .

RUN /opt/conda/bin/conda config --set ssl_verify False
RUN pip install --upgrade pip --trusted-host pypi.org --trusted-host files.pythonhosted.org
RUN pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org torchvision
RUN pip install --trusted-host pypi.org --trusted-host files.pythonhosted.org tqdm
WORKDIR /workspace
RUN chmod -R a+w /workspace

ADD . /opt/pytorch_dist_mnist
ENTRYPOINT ["mpirun", "--allow-run-as-root", "-n", "4", "python", "-u","/opt/pytorch_dist_mnist/mpi_mnist_cpu.py"]

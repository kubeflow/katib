# This image will setup the container for compiling and running the ksched binary
# The flowlessly solver is already set up in the location needed by ksched

# Prepare base image
FROM ubuntu:14.04
MAINTAINER Yuji Oshima <yuji.oshima0x3fd@gmail.com>
RUN apt-get update
RUN apt-get --yes --force-yes install git-all build-essential cmake libgflags2 libgflags-dev libgoogle-glog0 libgoogle-glog-dev libboost-all-dev curl wget autoconf automake libtool curl make g++ unzip

# Set up the solver binary in the location "/usr/local/bin/flowlessly/" hardcoded in the placement/solver
RUN cd /root && git clone https://github.com/hasbro17/Flowlessly.git && cd /root/Flowlessly && cmake . && make
RUN mkdir -p /usr/local/bin/flowlessly/ && cp /root/Flowlessly/build/flow_scheduler /usr/local/bin/flowlessly/

# Enivronment variables for Go
ENV PATH $PATH:/usr/local/go/bin:/root/go-workspace/bin::/usr/local/bin/flowlessly/
ENV GOPATH /root/go-workspace

# Install Go and set up workspace for compiling binary
RUN cd /root && curl -O https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz && tar -xvf go1.9.linux-amd64.tar.gz && mv go /usr/local
RUN mkdir -p /root/go-workspace/src/github.com

# Set up kubernetes in go workspace
RUN mkdir /root/go-workspace/src/k8s.io && cd /root/go-workspace/src/github.com/ && git clone https://github.com/kubernetes/kubernetes.git
#&& mv /root/go-workspace/src/github.com/kubernetes /root/go-workspace/src/k8s.io/

# Set up proto stuff
RUN cd /root && mkdir proto-bin && cd proto-bin && wget "https://github.com/google/protobuf/releases/download/v3.0.0/protoc-3.0.0-linux-x86_64.zip" && unzip protoc-3.0.0-linux-x86_64.zip && cp bin/protoc /usr/local/bin/protoc

RUN : &&\
mkdir -p /root/go-workspace/src/k8s.io && \
cd /root/go-workspace/src/k8s.io && \
git clone https://github.com/kubernetes/apimachinery && \
git clone https://github.com/kubernetes/api && \
git clone https://github.com/kubernetes/client-go && \
git clone https://github.com/kubernetes/kube-openapi && \
go get github.com/go-openapi/spec && \
go get github.com/golang/glog && \
go get github.com/google/gofuzz && \
go get github.com/emicklei/go-restful && \
go get github.com/spf13/cobra && \
go get github.com/spf13/viper && \
go get github.com/labstack/echo && \
go get github.com/sirupsen/logrus && \
:

ADD . /root/go-workspace/src/github.com/osrg/dlk
RUN mkdir /root/go-workspace/src/github.com/osrg/dlk/bin
WORKDIR /root/go-workspace/src/github.com/osrg/dlk
RUN go build -o bin/dlkmanager ./dlkmanager
RUN go build -o bin/dlkctl ./dlkctl
ENV PATH $PATH:/root/go-workspace/src/github.com/osrg/dlk/bin

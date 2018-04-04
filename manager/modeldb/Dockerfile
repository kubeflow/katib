FROM node:latest

EXPOSE 3000

RUN apt-get update \
    && apt-get install -y \
        g++ \
        make \
        wget \
		git \
    && apt-get clean

WORKDIR /root

# Install Thrift 0.10.0
RUN wget -q http://archive.apache.org/dist/thrift/0.10.0/thrift-0.10.0.tar.gz && \
    tar -xzf thrift-0.10.0.tar.gz && \
    cd thrift-0.10.0 && \
    ./configure --without-python && \
    make && \
    ln -n ~/thrift-0.10.0/compiler/cpp/thrift /usr/local/bin/thrift

RUN cd / && git clone https://github.com/mitdbg/modeldb

ADD manager/modeldb/frontend /modeldb/frontend

WORKDIR /modeldb/frontend

RUN  npm install

RUN mkdir -p './thrift' && \
    thrift -r -out './thrift' -gen js:node '../thrift/ModelDB.thrift' && \
    npm install
ENTRYPOINT ["/modeldb/dockerbuild/wait_for_backend.sh"]
CMD ["backend"]

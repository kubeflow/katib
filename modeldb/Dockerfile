FROM thrift:0.11  AS build-env
RUN apt-get update && apt-get install -y git
RUN git clone https://github.com/mitdbg/modeldb
ADD modeldb/frontend /modeldb/frontend
WORKDIR /modeldb/frontend
RUN mkdir -p './thrift' && \
    thrift -r -out './thrift' -gen js:node '../thrift/ModelDB.thrift' 

FROM node:alpine
COPY --from=build-env /modeldb /modeldb
WORKDIR /modeldb/frontend
RUN npm install
EXPOSE 3000
ENTRYPOINT ["sh","/modeldb/dockerbuild/wait_for_backend.sh"]

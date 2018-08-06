docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --python_out=plugins=grpc:./python --go_out=plugins=grpc:. -I. api.proto
docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --grpc-gateway_out=logtostderr=true:. -I. api.proto
docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --swagger_out=logtostderr=true:. -I. api.proto

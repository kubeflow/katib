docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --go_out=plugins=grpc:. -I. api.proto

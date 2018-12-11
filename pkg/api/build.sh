for proto in api.proto health/health.proto; do
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --python_out=plugins=grpc:./python --go_out=plugins=grpc:. -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --plugin=protoc-gen-grpc=/usr/bin/grpc_python_plugin --python_out=./python --grpc_out=./python -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --grpc-gateway_out=logtostderr=true:. -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --swagger_out=logtostderr=true:. -I. $proto
done

set -x
set -e
for proto in api.proto; do
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --python_out=plugins=grpc:./python --go_out=plugins=grpc:. -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --plugin=protoc-gen-grpc=/usr/bin/grpc_python_plugin --python_out=./python --grpc_out=./python -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --grpc-gateway_out=logtostderr=true:. -I. $proto
  docker run -it --rm -v $PWD:$(pwd) -w $(pwd) znly/protoc --swagger_out=logtostderr=true:. -I. $proto
done
docker build -t protoc-gen-doc gen-doc/
docker run --rm -v $PWD/gen-doc:/out -v $PWD:/apiprotos protoc-gen-doc --doc_opt=markdown,api.md -I /protobuf -I /apiprotos api.proto
docker run --rm -v $PWD/gen-doc:/out -v $PWD:/apiprotos protoc-gen-doc --doc_opt=html,index.html -I /protobuf -I /apiprotos api.proto

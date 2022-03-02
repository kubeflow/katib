FROM pseudomuto/protoc-gen-doc
RUN apk add --no-cache curl && \
    mkdir -p /protobuf/google/protobuf && \
    for f in any duration descriptor empty struct timestamp wrappers; do \
      curl -L -o /protobuf/google/protobuf/${f}.proto https://raw.githubusercontent.com/google/protobuf/master/src/google/protobuf/${f}.proto; \
    done && \
    mkdir -p /protobuf/google/api && \
    for f in annotations http; do \
      curl -L -o /protobuf/google/api/${f}.proto https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/master/third_party/googleapis/google/api/${f}.proto; \
    done && \
    mkdir -p /protobuf/github.com/gogo/protobuf/gogoproto && \
    curl -L -o /protobuf/github.com/gogo/protobuf/gogoproto/gogo.proto https://raw.githubusercontent.com/gogo/protobuf/master/gogoproto/gogo.proto && \
    mkdir -p /protobuf/github.com/mwitkow/go-proto-validators && \
    curl -L -o /protobuf/github.com/mwitkow/go-proto-validators/validator.proto https://raw.githubusercontent.com/mwitkow/go-proto-validators/master/validator.proto && \
    apk del --purge curl

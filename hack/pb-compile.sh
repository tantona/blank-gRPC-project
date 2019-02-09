#! /bin/bash

for filepath in $(find proto -name '*.proto' -type f) ; do
  protoc \
    -I/usr/local/include \
    -I. \
    -I$GOPATH/src \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
    -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
    --go_out=plugins=grpc:. \
    --grpc-gateway_out=logtostderr=true:. \
    --swagger_out=logtostderr=true:. \
    ${filepath}

  echo "PROTOC: compiled $filepath"
done

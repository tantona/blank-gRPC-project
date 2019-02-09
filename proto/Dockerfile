FROM golang:1.10.3

# For watchin protobuf definitions
RUN apt-get update && apt-get install -y inotify-tools unzip

# Install protoc compiler for dev
RUN cd /tmp && \
  curl -sSL https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip > /tmp/protoc-3.6.1-linux-x86_64.zip \
  && unzip /tmp/protoc-3.6.1-linux-x86_64.zip \
  && mv /tmp/bin/protoc /usr/local/bin/protoc \
  && mv /tmp/include/* /usr/local/include


RUN curl -sL https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v1.5.1/protoc-gen-grpc-gateway-v1.5.1-linux-x86_64 -o /go/bin/protoc-gen-grpc-gateway \
  && chmod 0755 /go/bin/protoc-gen-grpc-gateway

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
  && go get -u github.com/golang/protobuf/protoc-gen-go

# This is relative to the root of this project
COPY ./ /go/src/PROJECT_ROOT

WORKDIR /go/src/PROJECT_ROOT/

CMD ["./scripts/watch-proto.sh"]

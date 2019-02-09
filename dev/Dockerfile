FROM golang:latest

COPY ./ /go/src/PROJECT_ROOT/
WORKDIR /go/src/PROJECT_ROOT/

RUN go get github.com/pilu/fresh

CMD ["fresh", "-c", "/go/src/PROJECT_ROOT/hack/fresh.conf"]

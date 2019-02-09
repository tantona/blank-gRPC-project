FROM golang:1.10.3-alpine3.8 as build

RUN apk add -U --no-cache ca-certificates git bash

COPY ./ /go/src/PROJECT_ROOT/
WORKDIR /go/src/PROJECT_ROOT/

RUN go build -o app PROJECT_ROOT/services/environment/cmd \
  && mv ./app /go/bin

FROM alpine:3.8

WORKDIR /opt/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/app /opt/

CMD ["/opt/app"]

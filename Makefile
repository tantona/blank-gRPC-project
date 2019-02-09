.PHONY: test proto build clean fix awscli

build: clean proto
	@dep ensure -v
	@scripts/build-images.sh

test: proto
	@go test -v ./...

proto:
	@docker build -t compile-protocol-buffers -f proto/Dockerfile .
	@docker run -v ${PWD}:/go/src/PROJECT_ROOT compile-protocol-buffers ./hack/pb-compile.sh

clean:
	@rm -f ${PWD}/proto/services/*/*.go ${PWD}/proto/services/*/*.json
	@rm -rf ${PWD}/services/*/tmp

fix: clean proto
	@rm -rf vendor Gopkg.lock
	@dep ensure -v

init:
	@./hack/init.sh

base_dir = $(shell pwd)
gopath = "$(base_dir)/vendor:$(GOPATH)"

all: check-gopath clean test fmt build

build:
	@echo "==> Compiling source code."
	@env GOPATH=$(gopath) go build -v -o ./bin/ec2_metadata_mock ./ec2_metadata_mock

binaries: check-gopath clean test
	@echo "==> Compiling source code."
	@env GOPATH=$(gopath) GOOS=linux go build -v -o ./bin/ec2_metadata_mock-linux ./ec2_metadata_mock
	@env GOPATH=$(gopath) GOOS=darwin go build -v -o ./bin/ec2_metadata_mock-darwin ./ec2_metadata_mock

test: check-gopath
	@echo "==> Running tests."
	@env GOPATH=$(gopath) go test -cover ./ec2_metadata_mock/...

deps: check-gopath
	@echo "==> Downloading dependencies."
	@env GOPATH=$(gopath) go get -d -v ./ec2_metadata_mock/...

	@echo "==> Removing .git and .bzr from vendor."
	@find ./vendor -type d -name .git | xargs rm -rf
	@find ./vendor -type d -name .bzr | xargs rm -rf
	@find ./vendor -type d -name .hg | xargs rm -rf

fmt:
	@echo "==> Formatting source code."
	@gofmt -w ./ec2_metadata_mock

clean:
	@echo "==> Cleaning up previous builds."
	@rm -rf bin/ec2_metadata_mock

help:
	@echo "clean\t\tremove previous builds"
	@echo "deps\t\tdownload dependencies"
	@echo "fmt\t\tformat the code"
	@echo "test\t\ttest the code"
	@echo ""
	@echo "default will test, format, and build the code"

check-gopath:
ifndef GOPATH
  $(error GOPATH is undefined)
endif

.PNONY: all clean deps fmt help test

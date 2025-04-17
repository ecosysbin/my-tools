#
# Copyright 2023 The Zetyun.GCP Authors.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

SHELL = bash

OS = $(shell uname -s)
ARCH = $(shell uname -m)

# Base Variables.
REPO = $(shell go list -m)
LOGLEVEL ?= debug

CONFIGPATH ?= ./config.yaml

# Variables of Images.
IMAGE_REPO ?= harbor.zetyun.cn/gcp
IMAGE_NAME ?= vcluster-gateway
IMAGE_TAG ?= v0.0.0

# Variables of go build args and docker build args.
VERSION = $(IMAGE_TAG)
GIT_COMMIT = $(shell git rev-parse --short=11 HEAD)
DATE=`date "+%Y-%m-%d %H:%M:%S"`

# Generate go build args.
GO_BUILD_ARGS = \
	-w \
	-X '$(REPO)/pkg/version.Version=$(VERSION)' \
	-X '$(REPO)/pkg/version.GitCommit=$(GIT_COMMIT)' \
	-X '$(REPO)/pkg/version.BuildAt=$(DATE)' \

# Generate docker build args.
BUILD_ARGS= \
	--build-arg GO_BUILD_ARGS="$(GO_BUILD_ARGS)" \

##@ Build

#@ The targets of Build build and push images.
#@ Tags IMAGE_REPO IMAGE_NAME IMAGE_TAG can be overwrite.
#@ Examples:
#@     Build image with tag 'IMAGE_TAG'.
#@     - make docker-build IMAGE_TAG=v0.0.1
#@     Build image with all tags.
#@     - make docker-build IMAGE_REPO=<image-repo> IMAGE_NAME=<image-name> IMAGE_TAG=v0.0.1
#@     Build image and push with all tags.
#@     - make build-push IMAGE_REPO=<image-repo> IMAGE_NAME=<image-name> IMAGE_TAG=v0.0.1
#@

.PHONY: build-push
build-push: docker-build docker-push ## Build and push image to the harbor.

.PHONY: docker-build
docker-build: ## Build the docker image.
	docker build $(BUILD_ARGS) -t $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: docker-push
docker-push: ## Push the docker image.
	docker push $(IMAGE_REPO)/$(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: bin-darwin
bin-darwin: ## Build a binary file for macOS.
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -mod=vendor -ldflags "$(GO_BUILD_ARGS)" -o build/bin/vcluster-gateway-gateway cmd/vcluster_gateway.go

.PHONY: bin-linux
bin-linux: ## Build a binary file for linux.
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "$(GO_BUILD_ARGS)" -o build/bin/vcluster-gateway-gateway cmd/vcluster_gateway.go

.PHONY: add-minio-bucket
add-minio-bucket:
	mc config host add gcp-vcluster-chart http://172.20.58.69:9000 x6r5tYUtXOr1fCODRfbR loEx8nRJtR6po53zr9iE7ez4Y95XzbDn71OfUzNQ --api s3v4

.PHONY: package-chart
package-chart: ## Build package vcluster-k8s chart.
	cd charts && helm package vcluster-k8s --version 0.18.1
	cd charts && helm repo index .

.PHONY: push-chart
push-chart: package-chart ## Push chart to minio.
	cd charts && mc cp ./index.yaml gcp-vcluster-chart/charts/vcluster/k8s/charts-dev/ --recursive
	cd charts && mc cp ./vcluster-k8s-0.18.1.tgz gcp-vcluster-chart/charts/vcluster/k8s/charts-dev/ --recursive

##@ Deployment

.PHONY: gen-protos
gen-protos: ## Generate grpc files with osm/vcluster protos.
	buf generate --template buf.gen.yaml https://gitlab.datacanvas.com/AlayaNeW/OSM/apis.git#branch=vcluster --path datacanvas/gcp/osm/vcluster_1.1/v1
	find ./pkg/apis/grpc/gen -type f -name "*.pb.go" -print0 | xargs -0 -I {} protoc-go-inject-tag -input="{}"

.PHONY: install-minio-mc
install-minio-mc: ## Install minio cli mc.
ifeq (${OS}, Darwin)
	echo "Detected macOS"
	brew install minio/stable/mc
else ifeq (${OS}, Linux)
	echo "Detected Linux"
	ifeq (${ARCH}, x86_64)
		echo "Detected 64-bit Linux (x86_64)"
		curl -L https://dl.min.io/client/mc/release/linux-amd64/mc --create-dirs -o $${HOME}/minio-binaries/mc
		chmod +x $${HOME}/minio-binaries/mc
		export PATH=$${PATH}:$${HOME}/minio-binaries
	else ifeq (${ARCH}, ppc64le)
		echo "Detected 64-bit Linux (ppc64le)"
		curl -L https://dl.min.io/client/mc/release/linux-ppc64le/mc --create-dirs -o $${HOME}/minio-binaries/mc
		chmod +x $${HOME}/minio-binaries/mc
		export PATH=$${PATH}:$${HOME}/minio-binaries
	else ifeq (${ARCH}, aarch64)
		echo "Detected 64-bit Linux (aarch64)"
		curl -L https://dl.min.io/client/mc/release/linux-arm64/mc --create-dirs -o $${HOME}/minio-binaries/mc
		chmod +x $${HOME}/minio-binaries/mc
		export PATH=$${PATH}:$${HOME}/minio-binaries
	endif
endif

.PHONY: install-swag
install-swag: ## Install swag tool.
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swag
swag: ## Format swag annotations and init or update swagger docs.
	swag fmt
	swag init -g cmd/vcluster_gateway.go

.PHONY: tidy
tidy: ## Update dependencies and packages in vendor.
	go mod tidy
	go mod vendor

.PHONY: run
run: ## Run tests locally.
	go run -mod=vendor cmd/vcluster_gateway.go start --config=$(CONFIGPATH) --log-level=$(LOGLEVEL)

.PHONY: install-gofumpt
install-gofumpt: ## Install gofumpt tool.
	go install mvdan.cc/gofumpt@latest

.PHONY: format
format: ## Format go code using gofumpt
	gofumpt -l -w .

.PHONY: install-imports
install-imports: ## install goimports tool.
	go install -v github.com/incu6us/goimports-reviser/v3@latest

.PHONY: imports
imports: ## Sort go imports.
	goimports-reviser --rm-unused -project-name gitlab.datacanvas.com  -format ./...

.PHONY: check
check:
	go build -o /dev/null cmd/vcluster_gateway.go

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Show the targets Usage which in the Makefile.
	@awk 'BEGIN {FS = ":.*##"; printf "\033[1mUsage:\033[0m\n  \033[3mmake\033[0m \033[3;36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1;4m%s\033[0m\n", substr($$0, 5) } /^#@/ { printf "    \033[3;32m%s\033[0m\n", substr($$0, 4) } ' $(MAKEFILE_LIST)
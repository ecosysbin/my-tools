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

#FROM harbor.zetyun.cn/gcp/golang-alpine3.19-tini:1.22.1 as builder
FROM harbor.zetyun.cn/gcp/golang-alpine3.19-tini:1.22.1 as builder
WORKDIR /workspace

ENV GOPROXY  https://goproxy.cn,direct

#ENV GOINSECURE gitlab.datacanvas.com,git@gitlab.datacanvas.com
#ENV GONOPROXY gitlab.datacanvas.com/*,git@gitlab.datacanvas.com/*
#ENV GONOSUMDB gitlab.datacanvas.com,git@gitlab.datacanvas.com
#ENV GOPRIVATE gitlab.datacanvas.com/*,git@gitlab.datacanvas.com/*

COPY ./go.mod ./

RUN go mod download -x

COPY . ./

ARG GO_BUILD_ARGS

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GO_BUILD_ARGS" -o app-gateway cmd/app.go

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM harbor.zetyun.cn/aidc/bash-tini:stable

COPY --from=builder /workspace/app-gateway .
ENTRYPOINT [ "/usr/bin/tini", "--", "./app-gateway", "start" ]
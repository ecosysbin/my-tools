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

FROM registry.cn-beijing.aliyuncs.com/mfimg/private:golang1.22 as builder

WORKDIR /workspace

ENV GOPROXY https://goproxy.cn,direct
ARG GO_BUILD_ARGS

RUN wget http://containercloud-mirror.xaidc.com:9000/containercloud/helm && \
    chmod +x helm && \
    mv helm /usr/local/bin/helm

COPY charts/vcluster-k8s charts/vcluster-k8s

#ENV GOINSECURE gitlab.datacanvas.com,git@gitlab.datacanvas.com
#ENV GONOPROXY gitlab.datacanvas.com/*,git@gitlab.datacanvas.com/*
#ENV GONOSUMDB gitlab.datacanvas.com,git@gitlab.datacanvas.com
#ENV GOPRIVATE gitlab.datacanvas.com/*,git@gitlab.datacanvas.com/*

COPY go.mod .

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "$GO_BUILD_ARGS" -o vcluster-gateway cmd/vcluster_gateway.go

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM FROM registry.cn-beijing.aliyuncs.com/mfimg/private:alpine3.16

COPY --from=builder /workspace/vcluster-gateway .
COPY --from=builder /workspace/hack ./hack
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm
COPY --from=builder /workspace/charts/vcluster-k8s ./charts/vcluster-k8s

ENTRYPOINT [ "/usr/bin/tini", "--", "./vcluster-gateway", "start" ]
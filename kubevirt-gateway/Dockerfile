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

FROM harbor.zetyun.cn/gcp/golang-alpine3.19-tini:1.21.5 as builder

WORKDIR /workspace

COPY ./go.mod ./

RUN go mod download -x

COPY . ./

ARG GO_BUILD_ARGS

RUN go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags "$GO_BUILD_ARGS" -o kubevirt-gateway cmd/kubevirt_gateway.go

# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

FROM harbor.zetyun.cn/gcp/ubuntu-tini:22.04

COPY --from=builder /workspace/kubevirt-gateway .

ENTRYPOINT [ "/usr/bin/tini", "--", "./kubevirt-gateway" ]
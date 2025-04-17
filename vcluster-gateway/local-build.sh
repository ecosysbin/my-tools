#!/bin/zsh
set -o errexit
set -o nounset

# 1. 自动从 Git 中获取版本号（tag 或 commit）
#    --tags：优先获取 tag；--always：如果没有 tag 就使用 commit hash。
#    --dirty：如果工作区有改动会加上 -dirty 标识
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "unknown")

REPO=$(go list -m)
GIT_COMMIT=$(git rev-parse --short=11 HEAD)
DATE=$(date "+%Y-%m-%d %H:%M:%S")

GO_BUILD_ARGS="-w \
    -X '${REPO}/pkg/version.Version=${VERSION}' \
    -X '${REPO}/pkg/version.GitCommit=${GIT_COMMIT}' \
    -X '${REPO}/pkg/version.BuildAt=${DATE}'"

# 清理旧的可执行文件
rm -rf ./vcluster-gateway

# 2. 构建 Go 项目
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${GO_BUILD_ARGS}" -o vcluster-gateway cmd/vcluster_gateway.go

# 3. 构建 Docker 镜像
docker build -f Dockerfile-local -t "vc-g:${VERSION}" .

# 4. 推送到 Harbor
docker tag "vc-g:${VERSION}" "harbor.zetyun.cn/gcp/vcluster-gateway-local-build:${VERSION}"
docker push "harbor.zetyun.cn/gcp/vcluster-gateway-local-build:${VERSION}"

# 复制镜像名到剪贴板
echo "harbor.zetyun.cn/gcp/vcluster-gateway-local-build:${VERSION}" | pbcopy

# 5. 输出提示信息
echo -e "\033[32m构建完成并推送 Docker 镜像到 harbor.zetyun.cn/gcp/vcluster-gateway-local-build:${VERSION}\033[0m"
echo -e "\033[32m$(date "+%Y-%m-%d %H:%M:%S")\033[0m"

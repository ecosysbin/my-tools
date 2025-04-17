PLATFORM = linux/amd64
VERSION = latest

# 支持通过 make docker_build version=xxx 的方式传递版本号
.PHONY: docker_build
docker_build:
	docker build --platform ${PLATFORM} -t harbor.zetyun.cn/aidc/vcluster/vcluster-plugin:$(if $(version),$(version),$(VERSION)) .

# 推送镜像到 Docker 仓库
.PHONY: docker_push
docker_push:
	docker push harbor.zetyun.cn/aidc/vcluster/vcluster-plugin:$(if $(version),$(version),$(VERSION))

.PHONY: install-imports
install-imports: ## install goimports tool.
	go install -v github.com/incu6us/goimports-reviser/v3@latest

.PHONY: imports
imports: ## Sort go imports.
	goimports-reviser --rm-unused -project-name gitlab.datacanvas.com  -format ./...
generate: ## 执行generate 命令：包括生成Mock文件, 自动化生成error code
	go generate ./...


BUMP_MINOR_VERSION          ?= "$(shell gsemver bump patch)"


.PHONY: minor_version
minor_version: ## make MINOR_VERSION
	git tag "${BUMP_MINOR_VERSION}"
	git push origin ${BUMP_MINOR_VERSION}

help: ## 展示帮助信息
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.32.2)

lint:
	@echo "Running $@ check"
	@GO111MODULE=on golangci-lint cache clean
	@GO111MODULE=on golangci-lint run --timeout=5m --config ./.golangci.yml --fix
verifiers: getdeps lint

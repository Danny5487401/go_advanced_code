generate: ## 执行generate 命令：包括生成Mock文件, 自动化生成error code
	go generate ./...


BUMP_MINOR_VERSION          ?= "$(shell gsemver bump patch)"


.PHONY: minor_version
minor_version: ## make MINOR_VERSION
	git tag "${BUMP_MINOR_VERSION}"
	git push origin ${BUMP_MINOR_VERSION}

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"


getdeps:
	@mkdir -p ${GOPATH}/bin
	@which golangci-lint 1>/dev/null || (echo "Installing golangci-lint" && go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.32.2)

lint:
	@echo "Running $@ check"
	@GO111MODULE=on golangci-lint cache clean
	@GO111MODULE=on golangci-lint run --timeout=5m --config ./.golangci.yml --fix
verifiers: getdeps lint


ROOT_PACKAGE=github.com/Danny5487401/go_advanced_code

# ==============================================================================
# Includes

include scripts/make-rules/common.mk
include scripts/make-rules/tools.mk
include scripts/make-rules/golang.mk



.PHONY: format
format: tools.verify.golines tools.verify.goimports
	@echo "===========> Formating codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(FIND) -type f -name '*.go' | $(XARGS) golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated .


.PHONY:doctoc
doctoc: # 给当前目录及子目录的所有文件添加目录
	doctoc .
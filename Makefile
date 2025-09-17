.DEFAULT_GOAL:=help
-include .makerc

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: .mise
	@:

.PHONY: .mise
# Install dependencies
.mise:
	@mise install -q

### Tasks

.PHONY: check
## Run tests and linters
check: tidy lint test

.PHONY: tidy
## Run go mod tidy
tidy:
	@go mod tidy

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@go list -u -m -json all | go-mod-outdated -update -direct

.PHONY: test
## Run tests
test:
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out --tags=safe -race ./...

.PHONY: test.demo
## Run tests
test.demo: install
	@rm -rf tmp/test
	@mkdir -p tmp/test
	@cd tmp/test && \
		git init . && \
		git remote add origin https://github.com/foomo/posh-test-demo && \
		posh init && \
		echo "replace github.com/foomo/posh => ../../../" >> .posh/go.mod && \
		make shell.build && \
		bin/posh execute welcome demo

.PHONY: lint
## Run linter
lint:
	@golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	@golangci-lint run --fix

.PHONY: lint.super
## Run super linter
lint.super:
	docker run --rm -it \
		-e 'RUN_LOCAL=true' \
		-e 'DEFAULT_BRANCH=main' \
		-e 'IGNORE_GITIGNORED_FILES=true' \
		-e 'VALIDATE_JSCPD=false' \
		-e 'VALIDATE_GO=false' \
		-v $(PWD):/tmp/lint \
		github/super-linter

.PHONY: install
## Run go install
install: GOPATH=${shell go env GOPATH}
install:
	@go install -a main.go
	@mv "${GOPATH}/bin/main" "${GOPATH}/bin/posh"

.PHONY: install.debug
## Run go install with debug
install.debug:
	@go install -a -gcflags "all=-N -l" main.go

### Utils

.PHONY: help
## Show help text
help:
	@echo "\033[1;36mPOSH - Project oriented shell\033[0m"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "\033[36m%-23s\033[0m %s\n\n", cmd, help; help=""; \
			printf "\n\033[1;36m%s\033[0m\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  \033[36m%-23s\033[0m %s\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        " help "\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)


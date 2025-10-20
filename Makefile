.DEFAULT_GOAL:=help
-include .makerc

# --- Config -----------------------------------------------------------------

# Newline hack for error output
define br


endef

# --- Targets -----------------------------------------------------------------

# This allows us to accept extra arguments
%: .mise .husky
	@:

.PHONY: .mise
# Install dependencies
.mise: msg := $(br)$(br)Please ensure you have 'mise' installed and activated!$(br)$(br)$$ brew update$(br)$$ brew install mise$(br)$(br)See the documentation: https://mise.jdx.dev/getting-started.html$(br)$(br)
.mise:
ifeq (, $(shell command -v mise))
	$(error ${msg})
endif
	@mise install

.PHONY: .husky
# Configure git hooks for husky
.husky:
	@git config core.hooksPath .husky

### Tasks

.PHONY: check
## Run lint & tests
check: tidy lint test test.demo

.PHONY: tidy
## Run go mod tidy
tidy:
	@echo "〉go mod tidy"
	@go mod tidy

.PHONY: lint
## Run linter
lint:
	@echo "〉golangci-lint run"
	@golangci-lint run

.PHONY: lint.fix
## Fix lint violations
lint.fix:
	@echo "〉golangci-lint run fix"
	@golangci-lint run --fix

.PHONY: test
## Run tests
test:
	@echo "〉go test"
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out --tags=safe -race ./...

.PHONY: test.demo
## Run demo tests
test.demo: install
	@echo "〉testing demo"
	@rm -rf tmp/test
	@mkdir -p tmp/test
	@cd tmp/test && \
		git init . && \
		git remote add origin https://github.com/foomo/posh-test-demo && \
		posh init && \
		echo "replace github.com/foomo/posh => ../../../" >> .posh/go.mod && \
		make shell.build && \
		bin/posh execute welcome demo

.PHONY: build
## Build binary
build:
	@echo "〉building bin/posh"
	@rm -f bin/posh
	@go build -o bin/posh main.go

.PHONY: build.debug
## Build binary in debug mode
build.debug:
	@echo "〉building debug bin/posh"
	@rm -f bin/posh
	@go build -gcflags "all=-N -l" -o bin/posh main.go

.PHONY: install
## Run go install
install: GOPATH=${shell go env GOPATH}
install:
	@echo "〉installing $$GOPATH/bin/posh"
	@go install -a main.go
	@mv "${GOPATH}/bin/main" "${GOPATH}/bin/posh"

.PHONY: install.debug
## Run go install with debug
install.debug:
	@echo "〉installing debug $$GOPATH/bin/posh"
	@go install -a -gcflags "all=-N -l" main.go

.PHONY: outdated
## Show outdated direct dependencies
outdated:
	@echo "〉go mod outdated"
	@go list -u -m -json all | go-mod-outdated -update -direct

### Utils

.PHONY: help
## Show help text
help:
	@echo "Project Oriented SHELL (posh)\n"
	@echo "Usage:\n  make [task]"
	@awk '{ \
		if($$0 ~ /^### /){ \
			if(help) printf "%-23s %s\n\n", cmd, help; help=""; \
			printf "\n%s:\n", substr($$0,5); \
		} else if($$0 ~ /^[a-zA-Z0-9._-]+:/){ \
			cmd = substr($$0, 1, index($$0, ":")-1); \
			if(help) printf "  %-23s %s\n", cmd, help; help=""; \
		} else if($$0 ~ /^##/){ \
			help = help ? help "\n                        " substr($$0,3) : substr($$0,3); \
		} else if(help){ \
			print "\n                        " help "\n"; help=""; \
		} \
	}' $(MAKEFILE_LIST)

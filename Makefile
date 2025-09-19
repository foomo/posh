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
## Run lint & test
check: tidy lint test test.demo

.PHONY: tidy
## Run go mod tidy
tidy:
	@go mod tidy

.PHONY: lint
## Run linter
lint:
	@golangci-lint run

.PHONY: lint.fix
## Run linter and fix
lint.fix:
	@golangci-lint run --fix

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

.PHONY: build
## Build binary
build:
	@rm -f bin/posh
	@go build -o bin/posh main.go

.PHONY: build.debug
## Build binary in debug mode
build.debug:
	@rm -f bin/posh
	@go build -gcflags "all=-N -l" -o bin/posh main.go

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

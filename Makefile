.DEFAULT_GOAL:=help

## === Tasks ===

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
	@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -race ./...
	#@GO_TEST_TAGS=-skip go test -coverprofile=coverage.out -race -json ./... | gotestfmt

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

## === Utils ===

## Show help text
help:
	@awk '{ \
			if ($$0 ~ /^.PHONY: [a-zA-Z\-\_0-9]+$$/) { \
				helpCommand = substr($$0, index($$0, ":") + 2); \
				if (helpMessage) { \
					printf "\033[36m%-23s\033[0m %s\n", \
						helpCommand, helpMessage; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^[a-zA-Z\-\_0-9.]+:/) { \
				helpCommand = substr($$0, 0, index($$0, ":")); \
				if (helpMessage) { \
					printf "\033[36m%-23s\033[0m %s\n", \
						helpCommand, helpMessage"\n"; \
					helpMessage = ""; \
				} \
			} else if ($$0 ~ /^##/) { \
				if (helpMessage) { \
					helpMessage = helpMessage"\n                        "substr($$0, 3); \
				} else { \
					helpMessage = substr($$0, 3); \
				} \
			} else { \
				if (helpMessage) { \
					print "\n                        "helpMessage"\n" \
				} \
				helpMessage = ""; \
			} \
		}' \
		$(MAKEFILE_LIST)

-include .makerc
.DEFAULT_GOAL:=help

# --- .makerc -----------------------------------------------------------------
#
# level=debug

# --- Config -----------------------------------------------------------------

level?=info

# --- Helpers -----------------------------------------------------------------

.PHONY: bin/posh
# Builds posh and takes the git hash to detect changes
bin/posh: current=$(shell bin/posh version)
bin/posh: version=$(shell git ls-files -s .posh/pkg | git hash-object --stdin)
bin/posh: commitHash=$(shell git rev-parse HEAD)
bin/posh: buildTimestamp=$(shell date +%s)
bin/posh: ldflags=\
	-X github.com/foomo/posh/internal/version.Version=${version} \
	-X github.com/foomo/posh/internal/version.CommitHash=${commitHash} \
	-X github.com/foomo/posh/internal/version.BuildTimestamp=${buildTimestamp}
bin/posh:
	@if [ "${current}" != "${version}" ]; then \
  	cd .posh && go mod tidy && go build -trimpath -ldflags="${ldflags}" -o ../bin/posh main.go; \
  fi

# --- Targets -----------------------------------------------------------------

.PHONY: clean
## Remove built targets
clean:
	@rm -f bin/*

.PHONY: config
## Print posh config
config: bin/posh
	@bin/posh config --level ${level}

.PHONY: brew
## Install project specific packages
brew: bin/posh
	@bin/posh brew --level ${level}

.PHONY: require
## Validate dependencies
require: bin/posh
	@bin/posh require --level ${level}

.PHONY: shell
## Start the interactive
shell: bin/posh require brew
	@bin/posh prompt --level ${level}

.PHONY: shell.rebuild
## Rebuild and start the interactive
shell.rebuild: clean shell

## === Utils ===

.PHONY: help
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

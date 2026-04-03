include makefiles/core.mk
include makefiles/testing.mk
include makefiles/production.mk
include makefiles/development.mk

.PHONY: all setup test release clean

all: build-production-linux

setup: initialize-environment

test: run-all-unit-tests

clean: purge-artifacts

release: \
	purge-artifacts \
	build-production-linux \
	build-production-windows \
	build-production-cli

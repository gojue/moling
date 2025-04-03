include variables.mk
include functions.mk

.PHONY: all | env
all: clean build
	@echo $(shell date)

.ONESHELL:
SHELL = /bin/bash

.PHONY: env
env:
	@echo ---------------------------------------
	@echo "MoLing Makefile Environment:"
	@echo ---------------------------------------
	@echo "SNAPSHOT_VERSION         $(SNAPSHOT_VERSION)"
	@echo ---------------------------------------
	@echo "OS_NAME                  $(OS_NAME)"
	@echo "OS_ARCH                  $(OS_ARCH)"
	@echo "TARGET_OS                $(TARGET_OS)"
	@echo "TARGET_ARCH              $(TARGET_ARCH)"
	@echo "GO_VERSION               $(GO_VERSION)"
	@echo ---------------------------------------
	@echo "CMD_GIT                  $(CMD_GIT)"
	@echo "CMD_GO                   $(CMD_GO)"
	@echo "CMD_INSTALL              $(CMD_INSTALL)"
	@echo "CMD_MD5                  $(CMD_MD5)"
	@echo ---------------------------------------
	@echo "VERSION_NUM              $(VERSION_NUM)"
	@echo "LAST_GIT_TAG             $(LAST_GIT_TAG)"
	@echo ---------------------------------------


.PHONY: help
help:
	@echo "# environment"
	@echo "    $$ make env					# show makefile environment/variables"
	@echo ""
	@echo "# build"
	@echo "    $$ make all					# build MoLing"
	@echo ""
	@echo "# clean"
	@echo "    $$ make clean				# wipe ./bin/"
	@echo ""
	@echo "# test"

.PHONY: clean build

.PHONY: clean
clean:
	$(CMD_RM) -f $(OUT_BIN)*

.PHONY: build
build:clean
	$(call gobuild,$(TARGET_OS),$(TARGET_ARCH))

# Format the code
.PHONY: format
format:
	@echo "  ->  Formatting code"
	golangci-lint run --disable-all -E errcheck -E staticcheck

.PHONY: test
test:
	CGO_ENABLED=1 go test -v -race ./...

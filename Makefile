.PHONY: init configure add-commit-hook dev build-release test test_cov lint \
	pb_pkg_gen pb_internal_gen pb_model_gen pb_service_gen

PKG_PATH=go.neonyx.io
PROJECT=go-swn

VERSION = 0
PATCHLEVEL = 0
SUBLEVEL = 1
EXTRAVERSION = -rc1

DEV_VERSION = $(VERSION)$(if $(PATCHLEVEL),.$(PATCHLEVEL)$(if $(SUBLEVEL),.$(SUBLEVEL)))$(EXTRAVERSION)
RELEASE_VERSION = $(VERSION)$(if $(PATCHLEVEL),.$(PATCHLEVEL)$(if $(SUBLEVEL),.$(SUBLEVEL)))

DEPLOYMENT_DIR=./deployment

init: configure add-commit-hook
	@if [ -f go.mod ]; then \
		echo "[!] go.mod exists"; \
	else \
		go mod init $(PKG_PATH)/$(PROJECT); \
	fi

configure:
	@bash $(DEPLOYMENT_DIR)/configure.sh

add-commit-hook:
	@echo	"[*] Installing git commit hook..."
	pre-commit install -t commit-msg

dev:
	@echo	"[*] Building Docker image for development..."
	docker build -t $(PROJECT):$(DEV_VERSION) .

	@echo	"[*] Running Docker container: $(PROJECT)..."
	docker run --rm --name $(PROJECT) $(ARGS) $(PROJECT):$(DEV_VERSION)

build-release:
	@echo	"[*] Building Docker image for release..."
	docker build --build-arg GO_FLAGS=-ldflags="-s -w" -t $(PROJECT):$(RELEASE_VERSION) .

test:
	@echo	"[*] Running Go tests..."
	go test ./... $(ARGS)

test_cov:
	@echo	"[*] Running Go tests with coverage..."
	go test ./... -coverprofile=cover.out
	go tool cover -html=cover.out

lint:
	@echo	"[*] Running linters..."
	golangci-lint run

## Protobuf targets ##
include $(DEPLOYMENT_DIR)/proto.mk

# generate all pkg/* models and services .proto files
pb_pkg_gen: pb_pkg_model_gen pb_pkg_service_gen

# generate all internal/* models and services .proto files
pb_internal_gen: pb_internal_model_gen pb_internal_service_gen

# generate a given model protobuf in PROTO= CLI argument
pb_model_gen:
	$(call protoc_model,$(PROTO))

# generate a given service protobuf in PROTO= CLI argument
pb_service_gen:
	$(call protoc_service,$(PROTO))
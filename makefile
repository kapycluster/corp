# Makefile

# Directories
KAPYSERVER_DIR := ./kapyserver
PANEL_DIR := ./panel
CONTROLLER_DIR := ./controller

# Binaries
KAPYSERVER_BIN := bin/kapyserver
PANEL_BIN := bin/panel
CONTROLLER_BIN := bin/controller
CONTROLLER_GEN_BIN := bin/controller-gen

# Go settings
GO := go
GOBUILD := $(GO) build

# Kubebuilder settings
CONTROLLER_GEN := $(CONTROLLER_GEN_BIN)

# Default target
.PHONY: all
all: build

# Build all binaries
.PHONY: build
build: kapyserver panel controller

# Build kapyserver binary
.PHONY: kapyserver
kapyserver: $(KAPYSERVER_BIN)

$(KAPYSERVER_BIN):
	@echo "building kapyserver..."
	@mkdir -p bin
	$(GOBUILD) -o $@ $(KAPYSERVER_DIR)/cmd/main.go

# Build panel binary
.PHONY: panel
panel: $(PANEL_BIN)

$(PANEL_BIN):
	@echo "building panel..."
	@mkdir -p bin
	$(GOBUILD) -o $@ $(PANEL_DIR)/cmd/main.go

# Build controller binary
.PHONY: controller
controller: $(CONTROLLER_BIN)

$(CONTROLLER_BIN):
	@echo "building controller..."
	@mkdir -p bin
	$(GOBUILD) -o $@ $(CONTROLLER_DIR)/cmd/main.go

# Install controller-gen binary
.PHONY: install-controller-gen
install-controller-gen: $(CONTROLLER_GEN_BIN)

$(CONTROLLER_GEN_BIN):
	@echo "installing controller-gen..."
	@mkdir -p bin
	GOBIN=$(CURDIR)/bin $(GO) install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

# Generate Kubernetes manifests and types for controller
.PHONY: controller-gen
controller-gen: install-controller-gen
	@echo "controller: generating k8s manifests..."
	$(CONTROLLER_GEN) object paths="$(CONTROLLER_DIR)/..."
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd webhook paths="$(CONTROLLER_DIR)/..." output:crd:artifacts:config=$(CONTROLLER_DIR)/config/crd/bases

# Clean up binaries
.PHONY: clean
clean:
	@echo "cleaning up..."
	@rm -rf bin

# corp
> Monorepo for all Kapy Cluster code.

## overview

This repository contains code for 3 main components of Kapy Cluster, along
with various shared types and utilities.

1. `panel`: The dashboard monolith at panel.kapycluster.com. Mostly in Go,
  with templating in Templ and HTMX/Alpine.js on the frontend.
2. `controller`: The K8s controller that manages the `controlplanes.kapy.sh`
  CRD. Results in a `kapyserver` deployment in the management cluster.
3. `kapyserver`: The control plane server. Wrapper around K3s with a gRPC API
  for the controller to interact with. `kapyclient` provides the gRPC client.


## make targets

### default target
- `make all`: Builds all binaries (`kapyserver`, `panel`, and `controller`).
- `make clean`: Cleans up all built binaries.

### build targets
- `make build`: Builds all binaries (`kapyserver`, `panel`, and `controller`).
- `make kapyserver`: Builds the `kapyserver` binary.
- `make panel`: Builds the `panel` binary.
- `make controller`: Builds the `controller` binary.

### controller CRD generation and deployment
- `make install-controller-gen`: Installs the `controller-gen` binary.
- `make controller-gen`: Generates Kubernetes manifests and types for the controller.
- `make controller-install`: Installs the generated CRDs into the Kubernetes cluster.
- `make controller-uninstall`: Uninstalls the generated CRDs from the Kubernetes cluster.
- `make controller-deploy`: Deploys the controller to the Kubernetes cluster.
- `make controller-undeploy`: Undeploys the controller from the Kubernetes cluster.

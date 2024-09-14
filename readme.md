# corpy

All corporate kapy code lives here.


## overview

This is the up to date list (best effort) of the top-level directories
and their use:

- `controller`: Kubernetes controller that reconciles `ControlPlanes` (controlplanes.kapy.sh).
- `panel`: API layer and web frontend monolith.
- `kapyserver`: Kubernetes server. This is what they pay us for.
- `log`: Common logging package; wraps log/slog.
- `kapyclient`: gRPC client for kapyserver.
- `types`: Protobuf definitions for kapyserver's gRPC API. Could potentially contain more common types in the future.
- `docker`: Dockerfiles for the controller and kapyserver.

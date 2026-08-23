# Distro Crash Demo Service

A test service designed to simulate application issues with Odigos no-restart instrumentation for testing auto-rollback functionality.

## Behavior

1. **Initial State**: Starts successfully, serves HTTP requests on port 3000 (Node.js 20.19.2)
2. **Startup Check**: Checks for `ODIGOS_DISTRO_NAME` environment variable at startup
3. **Crash on Instrumentation**: If `ODIGOS_DISTRO_NAME` is set, the service crashes after an optional delay
4. **Delayed Crash**: Use `CRASH_DELAY_SECONDS` to crash after the service has been running (defaults to `0` for immediate crash at startup)
5. **Rollback Testing**: Designed to trigger Odigos automatic rollback after grace period

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ODIGOS_DISTRO_NAME` | No | - | When set, indicates Odigos instrumentation is active and triggers a crash |
| `CRASH_DELAY_SECONDS` | No | `0` | Seconds to wait before crashing. Only applies when `ODIGOS_DISTRO_NAME` is set |

## Currently built & pushed manually.

```bash
# Navigate to the service directory
cd tests/common/services/distro-crash-demo

# Build for AMD64 (GitHub Actions compatibility)
docker build --platform linux/amd64 -t distro-crash-demo:v1.0.0 .

# Alternative: Use buildx for multi-platform build
# docker buildx build --platform linux/amd64,linux/arm64 -t distro-crash-demo:v1.0.0 .

# Tag for ghcr.io
docker tag distro-crash-demo:v1.0.0 ghcr.io/odigos-io/simple-demo/odigos-demo-distro-crash:v1.0.0

# Push to GitHub Container Registry
docker push ghcr.io/odigos-io/simple-demo/odigos-demo-distro-crash:v1.0.0
```

## Kubernetes Manifests

All resources are deployed to the `distro-crash-demo` namespace.

Each variant has its own manifest file under `manifests/`:

| File | `CRASH_DELAY_SECONDS` |
|------|------------------------|
| `manifests/namespace.yaml` | dedicated namespace |
| `manifests/distro-crash-demo.yaml` | not set (immediate crash when instrumented) |
| `manifests/distro-crash-demo-10s.yaml` | `10` |
| `manifests/distro-crash-demo-60s.yaml` | `60` |
| `manifests/distro-crash-demo-120s.yaml` | `120` |

Each deployment file contains a Deployment, Service, and Odigos Source to instrument that deployment.

Apply a single variant:

```bash
kubectl apply -f manifests/namespace.yaml -f manifests/distro-crash-demo-60s.yaml
```

Apply all variants:

```bash
kubectl apply -f manifests/
```

Each deployment exposes HTTP probes on:
- `/healthz` — liveness and startup
- `/readyz` — readiness

## Testing Locally

```bash
# Run without instrumentation (should work fine)
docker run -p 3000:3000 distro-crash-demo:v1.0.0

# Run with simulated instrumentation (should crash immediately at startup)
docker run -p 3000:3000 -e ODIGOS_DISTRO_NAME=nodejs distro-crash-demo:v1.0.0

# Run with simulated instrumentation and delayed crash (serves requests first)
docker run -p 3000:3000 -e ODIGOS_DISTRO_NAME=nodejs -e CRASH_DELAY_SECONDS=30 distro-crash-demo:v1.0.0

# Test endpoints (works without instrumentation, or before delayed crash)
curl http://localhost:3000
```

## Usage in Tests

Use this service to test no-restart instrumentation rollback scenarios:

- **Immediate crash**: default behavior when instrumented (`CRASH_DELAY_SECONDS=0`)
- **Delayed crash**: pod stays healthy during grace period, then crashes after stability window scenarios

The service uses `ODIGOS_DISTRO_NAME` instead of `OTEL_SERVICE_NAME`, matching how Odigos marks instrumented workloads.

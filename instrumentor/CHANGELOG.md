# Changelog

## [1.1.0](https://github.com/blumamir/odigos/compare/v1.0.54...v1.1.0) (2024-04-11)


### Features

* add own telemetry pipeline deployment ([#641](https://github.com/blumamir/odigos/issues/641)) ([7bfcf92](https://github.com/blumamir/odigos/commit/7bfcf92eb309afa0b23441cd505f5c6540263a8f))
* configuring collector processors via crds ([#978](https://github.com/blumamir/odigos/issues/978)) ([4869bee](https://github.com/blumamir/odigos/commit/4869bee949d504ad95c60aa7f48fceee97037f4d))
* generic support for otel sdk types via crd ([#818](https://github.com/blumamir/odigos/issues/818)) ([461cdde](https://github.com/blumamir/odigos/commit/461cdde83ae2a4892a6a85a7521cc5a9df76320f))
* **odiglet:** support workloads for ebpf director ([#786](https://github.com/blumamir/odigos/issues/786)) ([431cbee](https://github.com/blumamir/odigos/commit/431cbee04af2e5b9a9fdd11995bd0a45fae48d25))
* reported name for source ([#369](https://github.com/blumamir/odigos/issues/369)) ([c731c91](https://github.com/blumamir/odigos/commit/c731c916ff5cc6cb5ddd90d54e556b4e92e20001))
* use annotation for ebpf instrumentation ([#766](https://github.com/blumamir/odigos/issues/766)) ([dd32c39](https://github.com/blumamir/odigos/commit/dd32c390512928c105b32b919506854728fb4291))
* use device manager for ebpf instrumentations ([#826](https://github.com/blumamir/odigos/issues/826)) ([fc3f417](https://github.com/blumamir/odigos/commit/fc3f41777cc42679ad9df60c25a9ba5db46875ce))


### Bug Fixes

* **instrumentor:** preserve containers with missing runtime details ([#844](https://github.com/blumamir/odigos/issues/844)) ([95340b5](https://github.com/blumamir/odigos/commit/95340b509ae4b9a07e745de31f9e9cdff572cddd))
* **instrumentor:** update instrumentation device on config default sdk config change ([#866](https://github.com/blumamir/odigos/issues/866)) ([0ee403f](https://github.com/blumamir/odigos/commit/0ee403f87e3c5a030b624d0698524116d1aa4f84))
* kind is compared with wrong case ([#805](https://github.com/blumamir/odigos/issues/805)) ([247d344](https://github.com/blumamir/odigos/commit/247d344ac4afbaaeb35ba969013dd8bf49638933))
* remove old instrumentation device when upgrading from community to enterprise ([#862](https://github.com/blumamir/odigos/issues/862)) ([509d175](https://github.com/blumamir/odigos/commit/509d17582b1189961e6f4e0ceb1d5fcaecfd024b))

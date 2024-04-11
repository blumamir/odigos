# Changelog

## [1.1.0](https://github.com/blumamir/odigos/compare/v1.0.54...v1.1.0) (2024-04-11)


### Features

* add memory limiter to gateway collector ([#1086](https://github.com/blumamir/odigos/issues/1086)) ([53591fc](https://github.com/blumamir/odigos/commit/53591fc3a2319325f4bf08f26da1c709e2a69297))
* **api:** instrumentation config crd ([#781](https://github.com/blumamir/odigos/issues/781)) ([10588c4](https://github.com/blumamir/odigos/commit/10588c4bdb4c7bba989b0b2615e5182104fb7b99))
* configuring collector processors via crds ([#978](https://github.com/blumamir/odigos/issues/978)) ([4869bee](https://github.com/blumamir/odigos/commit/4869bee949d504ad95c60aa7f48fceee97037f4d))
* frontend api for insertclusterattribute ([#1009](https://github.com/blumamir/odigos/issues/1009)) ([5fdd31d](https://github.com/blumamir/odigos/commit/5fdd31dcb8f3cccb0e35d729b061832a159518ab))
* generic support for otel sdk types via crd ([#818](https://github.com/blumamir/odigos/issues/818)) ([461cdde](https://github.com/blumamir/odigos/commit/461cdde83ae2a4892a6a85a7521cc5a9df76320f))
* implement rename action ([#1094](https://github.com/blumamir/odigos/issues/1094)) ([0bbc6a7](https://github.com/blumamir/odigos/commit/0bbc6a7259b34220359e4447535c0f49c320c82f))
* new action - delete attribute ([#1055](https://github.com/blumamir/odigos/issues/1055)) ([e12b524](https://github.com/blumamir/odigos/commit/e12b52431936dbae166c81c90ccbbfb1c793d337))
* **odiglet:** support workloads for ebpf director ([#786](https://github.com/blumamir/odigos/issues/786)) ([431cbee](https://github.com/blumamir/odigos/commit/431cbee04af2e5b9a9fdd11995bd0a45fae48d25))
* odigos action insert cluster attributes ([#1000](https://github.com/blumamir/odigos/issues/1000)) ([35da47b](https://github.com/blumamir/odigos/commit/35da47b13cd06de4c3958cdf6bd42394c6f7e60e))
* odigos migration framework ([#692](https://github.com/blumamir/odigos/issues/692)) ([3e0b3fd](https://github.com/blumamir/odigos/commit/3e0b3fd1c8cdcf451b548fab2f2976e569e253e7))
* store unique dest by id instead of user custom name ([4149981](https://github.com/blumamir/odigos/commit/41499811cf35f1366078d53a001f33572cb65d5e))


### Bug Fixes

* fix autoscaler reconciliation loop with patching collector groups instead of updating them ([#1089](https://github.com/blumamir/odigos/issues/1089)) ([834dc14](https://github.com/blumamir/odigos/commit/834dc1496e3d62f5b34c9ad8201cb8ff3a933e1a))
* processor CRD collector role enum new values and make orderHint optional ([#983](https://github.com/blumamir/odigos/issues/983)) ([f03620b](https://github.com/blumamir/odigos/commit/f03620b0d3e163dcb8493d9af774833e2ee9f6be))

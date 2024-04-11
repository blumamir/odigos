# Changelog

## [1.1.0](https://github.com/blumamir/odigos/compare/v1.0.54...v1.1.0) (2024-04-11)


### Features

* add k8s.pod.name resource attribute to go-ebpf instrumented pods ([#1033](https://github.com/blumamir/odigos/issues/1033)) ([64be5a9](https://github.com/blumamir/odigos/commit/64be5a986e1d6be28a87237668d9e7842f985c50))
* add own telemetry pipeline deployment ([#641](https://github.com/blumamir/odigos/issues/641)) ([7bfcf92](https://github.com/blumamir/odigos/commit/7bfcf92eb309afa0b23441cd505f5c6540263a8f))
* add workload resource attributes for native otel sdks ([#887](https://github.com/blumamir/odigos/issues/887)) ([4accebc](https://github.com/blumamir/odigos/commit/4accebc2aeb1259b7213b52dbee4f3221f5f108d))
* allow adding extra reconcilers to odiglet ([#785](https://github.com/blumamir/odigos/issues/785)) ([2fb4fa8](https://github.com/blumamir/odigos/commit/2fb4fa8707ccaab3c3582181abdfcc3074848c5c))
* **api:** instrumentation config crd ([#781](https://github.com/blumamir/odigos/issues/781)) ([10588c4](https://github.com/blumamir/odigos/commit/10588c4bdb4c7bba989b0b2615e5182104fb7b99))
* generic support for otel sdk types via crd ([#818](https://github.com/blumamir/odigos/issues/818)) ([461cdde](https://github.com/blumamir/odigos/commit/461cdde83ae2a4892a6a85a7521cc5a9df76320f))
* **odiglet:** support workloads for ebpf director ([#786](https://github.com/blumamir/odigos/issues/786)) ([431cbee](https://github.com/blumamir/odigos/commit/431cbee04af2e5b9a9fdd11995bd0a45fae48d25))
* use annotation for ebpf instrumentation ([#766](https://github.com/blumamir/odigos/issues/766)) ([dd32c39](https://github.com/blumamir/odigos/commit/dd32c390512928c105b32b919506854728fb4291))
* use device manager for ebpf instrumentations ([#826](https://github.com/blumamir/odigos/issues/826)) ([fc3f417](https://github.com/blumamir/odigos/commit/fc3f41777cc42679ad9df60c25a9ba5db46875ce))


### Bug Fixes

* add read permission for javaagent.jar ([#387](https://github.com/blumamir/odigos/issues/387)) ([540cc79](https://github.com/blumamir/odigos/commit/540cc79e6907c22e323839275efa0f8dcc1592d5)), closes [#385](https://github.com/blumamir/odigos/issues/385)
* broken context propagation on minikube ([#349](https://github.com/blumamir/odigos/issues/349)) ([8c08d24](https://github.com/blumamir/odigos/commit/8c08d240dd52e441c6086afd930726c83d986f2b))
* odiglet bind address already in use ([#531](https://github.com/blumamir/odigos/issues/531)) ([7ea83ee](https://github.com/blumamir/odigos/commit/7ea83ee7f59c9d2de426c5e586180f4c95bd45d0))
* **odiglet:** odiglet might use old instrumentation agents after version upgrade ([#395](https://github.com/blumamir/odigos/issues/395)) ([2fe8dfe](https://github.com/blumamir/odigos/commit/2fe8dfe9399a86295475cc4362fad65052758354))
* **odiglet:** only include non-terminating pods for runtime details analysis ([#870](https://github.com/blumamir/odigos/issues/870)) ([19a83d1](https://github.com/blumamir/odigos/commit/19a83d159fc97b73bbd1eab65d09a65a1b40e61c))
* pod reconciler for odiglet restart ([#863](https://github.com/blumamir/odigos/issues/863)) ([3232946](https://github.com/blumamir/odigos/commit/32329461a5a1ecafa1ed037fc5f89e61bf4ce629))

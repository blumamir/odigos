# Changelog

## [1.1.0](https://github.com/blumamir/odigos/compare/v1.0.54...v1.1.0) (2024-04-11)


### Features

* add aws and gcp resource attributes detector to data-collection collector ([#967](https://github.com/blumamir/odigos/issues/967)) ([0af6f78](https://github.com/blumamir/odigos/commit/0af6f78e88c81c58dbbb43270fcd8f9382f36497))
* add config option for loki labels ([#955](https://github.com/blumamir/odigos/issues/955)) ([58a3323](https://github.com/blumamir/odigos/commit/58a3323cf516297b663c9b04fd97607cacccc7f7))
* add generic otlp destination to odigos ui ([#1035](https://github.com/blumamir/odigos/issues/1035)) ([cc1d2fe](https://github.com/blumamir/odigos/commit/cc1d2fe5c71045bbb71b237de22124791cc4103b))
* add memory limiter to gateway collector ([#1086](https://github.com/blumamir/odigos/issues/1086)) ([53591fc](https://github.com/blumamir/odigos/commit/53591fc3a2319325f4bf08f26da1c709e2a69297))
* add otlp http destination with basic auth ([#1038](https://github.com/blumamir/odigos/issues/1038)) ([b7af0ee](https://github.com/blumamir/odigos/commit/b7af0ee0bfd7652420eb2ef39656d7889c889137))
* add own telemetry pipeline deployment ([#641](https://github.com/blumamir/odigos/issues/641)) ([7bfcf92](https://github.com/blumamir/odigos/commit/7bfcf92eb309afa0b23441cd505f5c6540263a8f))
* add workload resource attributes for native otel sdks ([#887](https://github.com/blumamir/odigos/issues/887)) ([4accebc](https://github.com/blumamir/odigos/commit/4accebc2aeb1259b7213b52dbee4f3221f5f108d))
* allow configuring loki labels for self managed loki destination ([#962](https://github.com/blumamir/odigos/issues/962)) ([8aa7ace](https://github.com/blumamir/odigos/commit/8aa7ace78ebc3351ec8095e9e9abf7a5c899a72a))
* **autoscaler:** Sync Data-Collection Affinity, Tolerations, and NodeSelector Attributes with Odiglet Daemonset ([#969](https://github.com/blumamir/odigos/issues/969)) ([91f15e5](https://github.com/blumamir/odigos/commit/91f15e51e37f7ff56c3e3f93b51b366280186ffe))
* configuring collector processors via crds ([#978](https://github.com/blumamir/odigos/issues/978)) ([4869bee](https://github.com/blumamir/odigos/commit/4869bee949d504ad95c60aa7f48fceee97037f4d))
* elasticsearch destination https and basic auth support ([#1042](https://github.com/blumamir/odigos/issues/1042)) ([a036000](https://github.com/blumamir/odigos/commit/a03600040982a9969c8a92b3a5170fd7d85a29c9))
* **elasticsearch:** remove hard-coded port 9200, validate URL input ([#941](https://github.com/blumamir/odigos/issues/941)) ([3b72dd1](https://github.com/blumamir/odigos/commit/3b72dd12e6f3c924bfdc7af9aa533347fd4a767a))
* frontend api for insertclusterattribute ([#1009](https://github.com/blumamir/odigos/issues/1009)) ([5fdd31d](https://github.com/blumamir/odigos/commit/5fdd31dcb8f3cccb0e35d729b061832a159518ab))
* generic support for otel sdk types via crd ([#818](https://github.com/blumamir/odigos/issues/818)) ([461cdde](https://github.com/blumamir/odigos/commit/461cdde83ae2a4892a6a85a7521cc5a9df76320f))
* grafanacloudloki ([#918](https://github.com/blumamir/odigos/issues/918)) ([64f0a7f](https://github.com/blumamir/odigos/commit/64f0a7feef9dd7c4b9d597a7d22d2a1443c52781))
* implement rename action ([#1094](https://github.com/blumamir/odigos/issues/1094)) ([0bbc6a7](https://github.com/blumamir/odigos/commit/0bbc6a7259b34220359e4447535c0f49c320c82f))
* new action - delete attribute ([#1055](https://github.com/blumamir/odigos/issues/1055)) ([e12b524](https://github.com/blumamir/odigos/commit/e12b52431936dbae166c81c90ccbbfb1c793d337))
* odigos action insert cluster attributes ([#1000](https://github.com/blumamir/odigos/issues/1000)) ([35da47b](https://github.com/blumamir/odigos/commit/35da47b13cd06de4c3958cdf6bd42394c6f7e60e))
* **quickwit:** quickwit destination ([#907](https://github.com/blumamir/odigos/issues/907)) ([2d55edd](https://github.com/blumamir/odigos/commit/2d55eddc76b64068cdf3ef482fa7b9eafc2fc5cb))
* use aws s3 exporter from contrib ([#966](https://github.com/blumamir/odigos/issues/966)) ([ba16a79](https://github.com/blumamir/odigos/commit/ba16a796932d565ba0224a0eb141399875af711e))


### Bug Fixes

* allow larger gRPC batches in otlp receiver ([#1048](https://github.com/blumamir/odigos/issues/1048)) ([802b332](https://github.com/blumamir/odigos/commit/802b3322f17a092d48e4307761f0ea7cffc8519f))
* allow non-standard port for Jaeger destination ([#809](https://github.com/blumamir/odigos/issues/809)) ([ef308a4](https://github.com/blumamir/odigos/commit/ef308a4f2b809b6f27d44f8b713d375090bb14c3))
* allow path in otlp http destinations ([#1045](https://github.com/blumamir/odigos/issues/1045)) ([24d9c08](https://github.com/blumamir/odigos/commit/24d9c0846c03a27ce3ac271e9f6c4e697dbd32d5))
* append destination processors after crd processors ([#988](https://github.com/blumamir/odigos/issues/988)) ([877eb00](https://github.com/blumamir/odigos/commit/877eb00af388cc79283b10ccc646b864eac7f584))
* **collector:** fail silently if kubelet access is not available ([#883](https://github.com/blumamir/odigos/issues/883)) ([bfafb1f](https://github.com/blumamir/odigos/commit/bfafb1f0fcb379338a16524fdd4e789aa4dca3fc))
* configure filelog with only instrumented apps' logs ([#1080](https://github.com/blumamir/odigos/issues/1080)) ([89eb546](https://github.com/blumamir/odigos/commit/89eb5466740131b2f701b4f8fdbfe1e49cc3fba4))
* fix autoscaler reconciliation loop with patching collector groups instead of updating them ([#1089](https://github.com/blumamir/odigos/issues/1089)) ([834dc14](https://github.com/blumamir/odigos/commit/834dc1496e3d62f5b34c9ad8201cb8ff3a933e1a))
* grafana cloud prometheus ([#925](https://github.com/blumamir/odigos/issues/925)) ([7d5b12f](https://github.com/blumamir/odigos/commit/7d5b12f826de14679a6dccb53ee8560060b8412b))
* grafana cloud tempo destination ([#917](https://github.com/blumamir/odigos/issues/917)) ([5228088](https://github.com/blumamir/odigos/commit/5228088773e076df5b3678aee487dfc55c210679))
* loki configuration  ([#879](https://github.com/blumamir/odigos/issues/879)) ([67e92ac](https://github.com/blumamir/odigos/commit/67e92acf284d97a7a653cd41cfbb0e0a708f3c1a))
* missing kubeletstats from collector ([#848](https://github.com/blumamir/odigos/issues/848)) ([48d5485](https://github.com/blumamir/odigos/commit/48d54850f012a32d577137e3162320c1408b1770))
* multiple destinations of the same kind ([#906](https://github.com/blumamir/odigos/issues/906)) ([5f14090](https://github.com/blumamir/odigos/commit/5f14090fa9cab58855c87dab4c97b2cda246c3f1))
* one failed daemonset pod prevents other from updating ([#1101](https://github.com/blumamir/odigos/issues/1101)) ([37987d1](https://github.com/blumamir/odigos/commit/37987d1f9849785519592d66ab7ab70a1dfeff65))
* quickwit destination no duplicate processors ([#1072](https://github.com/blumamir/odigos/issues/1072)) ([799a782](https://github.com/blumamir/odigos/commit/799a782998d7c68cb24b1a1887d442cbb9e74a79))
* remove default port 4318 from otlp http when missing ([#1056](https://github.com/blumamir/odigos/issues/1056)) ([411f09d](https://github.com/blumamir/odigos/commit/411f09d6d935ed4e6ca5086833a7bc97c535ae56))
* skip verify of kubelet stats metrics receiver secured endpoint ([#913](https://github.com/blumamir/odigos/issues/913)) ([534139d](https://github.com/blumamir/odigos/commit/534139da5c0a599fbdfad6755cccc1b8ba21b3c8))

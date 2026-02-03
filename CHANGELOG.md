# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.0.0 (2026-02-03)


### Features

* **content:** implement builder/encoding and complete PRD-173 ([5368a68](https://github.com/jonwraymond/toolprotocol/commit/5368a68c32392fd131518c7f62ec7168d9a9de2e))
* **content:** implement core types (GREEN) ([cb91beb](https://github.com/jonwraymond/toolprotocol/commit/cb91beb5f7f27b354ccd54ea73e03724fb6c16ec))
* **discover:** add negotiation and complete PRD-172 (GREEN) ([f71694f](https://github.com/jonwraymond/toolprotocol/commit/f71694f85aefdb5500a4cbfddb7d3083c4104f33))
* **discover:** implement core types (GREEN) ([b573151](https://github.com/jonwraymond/toolprotocol/commit/b573151fb374df663f8b381da280408186139341))
* **discover:** implement MemoryDiscovery (GREEN) ([60525a8](https://github.com/jonwraymond/toolprotocol/commit/60525a88a9c7461554c4fb531fd5092f9d9b609a))
* **discover:** implement Service (GREEN) ([0eb1600](https://github.com/jonwraymond/toolprotocol/commit/0eb16007e0263cb4063d4065fe95fc42ddca8dc0))
* **elicit:** implement PRD-177 user input elicitation package ([3c2deef](https://github.com/jonwraymond/toolprotocol/commit/3c2deef4838a3161323f0d2c913e6fef74934ee3))
* expose A2A package ([#7](https://github.com/jonwraymond/toolprotocol/issues/7)) ([f7cca67](https://github.com/jonwraymond/toolprotocol/commit/f7cca6708644c6c27c7efdbe217ba5115fa3dbe4))
* initial repository structure ([8ba2f15](https://github.com/jonwraymond/toolprotocol/commit/8ba2f156b9123d31a181262e25f96aa28155551e))
* **prompt:** implement PRD-179 MCP prompt template package ([844b7f3](https://github.com/jonwraymond/toolprotocol/commit/844b7f37060135b339c71d501fbdb1bb54b14410))
* **resource:** implement PRD-178 MCP resource package ([9837721](https://github.com/jonwraymond/toolprotocol/commit/9837721399905a38bbc8730e7988069bb0e7871d))
* **session:** implement PRD-176 session management package ([ce94b07](https://github.com/jonwraymond/toolprotocol/commit/ce94b07576992b1897230aa0fb96a8f7e11f5cfa))
* **stream:** implement PRD-175 streaming response package ([110a49e](https://github.com/jonwraymond/toolprotocol/commit/110a49e121ba03721dff4e596a11f68ec9384c3a))
* **task:** add errors and options, complete PRD-174 (GREEN) ([edb0669](https://github.com/jonwraymond/toolprotocol/commit/edb06693c15d855c23bcc362e33164ef0314c781))
* **task:** implement core types (GREEN) ([1764055](https://github.com/jonwraymond/toolprotocol/commit/17640551f98e3ef8e0a3394fd8331f1ac11f9eff))
* **task:** implement manager basic ops (GREEN) ([96c7e2d](https://github.com/jonwraymond/toolprotocol/commit/96c7e2d050fcd2e65e4639ffa839da1ff17544e2))
* **task:** implement memory store (GREEN) ([053cc98](https://github.com/jonwraymond/toolprotocol/commit/053cc987038bd14d08250ee56c57b4b80c008b32))
* **task:** implement state transitions (GREEN) ([418f888](https://github.com/jonwraymond/toolprotocol/commit/418f888121ccb9fc78ec993731077a8879c48b68))
* **task:** implement subscriptions (GREEN) ([bb0e953](https://github.com/jonwraymond/toolprotocol/commit/bb0e9539ef71a9620a7a90ba163f911ad316afc5))
* **toolprotocol:** add comprehensive TDD architecture improvements ([217b8dc](https://github.com/jonwraymond/toolprotocol/commit/217b8dc6995db663d15e37be7015939a3b53eeed))
* **transport:** add config types and errors (GREEN) ([75cbcc4](https://github.com/jonwraymond/toolprotocol/commit/75cbcc470e9ddd78cc09d8fdb28d553afa6f82d5))
* **transport:** add factory and registry (GREEN) ([746fe82](https://github.com/jonwraymond/toolprotocol/commit/746fe825ad8aa43a24c21b71fbe510585d20d3f3))
* **transport:** implement core types (GREEN) ([d1b8687](https://github.com/jonwraymond/toolprotocol/commit/d1b8687d3223e359cf5fc43b5065be1407cd3226))
* **transport:** implement SSE transport (GREEN) ([917c6ff](https://github.com/jonwraymond/toolprotocol/commit/917c6ffd8f3da921d35ea2ca87d2d9bab70ca9eb))
* **transport:** implement stdio transport (GREEN) ([88cfca2](https://github.com/jonwraymond/toolprotocol/commit/88cfca25ee150b3958e04dcc2796a6d01df5cea0))
* **transport:** implement streamable HTTP transport (GREEN) ([e974c17](https://github.com/jonwraymond/toolprotocol/commit/e974c170d94629cf9f757db09d646e2373e76c88))
* **wire:** add registry and errors, complete PRD-171 (GREEN) ([0488d28](https://github.com/jonwraymond/toolprotocol/commit/0488d288508ba6c4f56936c665e3f207d7d2ea01))
* **wire:** implement A2A wire format (GREEN) ([bc6dbf1](https://github.com/jonwraymond/toolprotocol/commit/bc6dbf125438afdd18d0d82cedc133f454cb6465))
* **wire:** implement ACP wire format (GREEN) ([a2f9f1d](https://github.com/jonwraymond/toolprotocol/commit/a2f9f1d14a3ed39850694e27d0d94c4d9258c98c))
* **wire:** implement core types (GREEN) ([f148926](https://github.com/jonwraymond/toolprotocol/commit/f14892666c4ebad9bd243b69db64ab2eb3499c58))
* **wire:** implement MCP wire format (GREEN) ([36f48ff](https://github.com/jonwraymond/toolprotocol/commit/36f48fff8c86da5803ec6dc24ea009e7f820010d))


### Bug Fixes

* **ci:** handle missing go.sum when no external dependencies ([71c1695](https://github.com/jonwraymond/toolprotocol/commit/71c1695e81ab1e17718c1da59705878c57395e3c))
* **lint:** address errcheck and gofmt issues ([7848aeb](https://github.com/jonwraymond/toolprotocol/commit/7848aeb6505671b9e1c6d7505878c6f4a0d34ec9))
* **stream:** make DefaultStream close race-safe ([7119d70](https://github.com/jonwraymond/toolprotocol/commit/7119d70b702497beba4845d66593b048ed8dc48e))


### Documentation

* add mkdocs config ([affe3d8](https://github.com/jonwraymond/toolprotocol/commit/affe3d84339f63e7b18b484f5e436468ef47dcbc))
* **toolprotocol:** align docs and mkdocs ([dd71d4f](https://github.com/jonwraymond/toolprotocol/commit/dd71d4f4c029989c49e5a17e3d1f0d38cf80e344))
* **toolprotocol:** align package docs and examples ([281587b](https://github.com/jonwraymond/toolprotocol/commit/281587bfd657a87bd7e56b0e91787f4309d3da28))
* update version matrix ([#9](https://github.com/jonwraymond/toolprotocol/issues/9)) ([75b9a1d](https://github.com/jonwraymond/toolprotocol/commit/75b9a1d4bbe72c0bba35def56390d54fd330c8a4))

## [Unreleased]

### Added
- Initial repository structure

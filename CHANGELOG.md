# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Update `architect-orb` to v4.8.0 to stop pushing to app collection helm chart.

## [2.2.0] - 2021-07-27

### Added

- Add `giantswarm.io/monitoring: "true"` label to enable Prometheus autodiscovery.

## [2.1.0] 2021-05-26

## [2.0.0] 2020-05-27

### Changed

- Update labels and selector labels.
- Change docker image templating.
- Change user/group ID templating.

### Removed

- Remove name and namespace from values.yaml.

## [1.0.0] 2019-08-13

## [0.1.0]  2019-08-13

[Unreleased]: https://github.com/giantswarm/vault-exporter/compare/v2.2.0...HEAD
[2.2.0]: https://github.com/giantswarm/vault-exporter/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/giantswarm/vault-exporter/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/giantswarm/vault-exporter/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/giantswarm/vault-exporter/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/giantswarm/vault-exporter/releases/tag/v0.1.0

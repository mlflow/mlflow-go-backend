# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1] - 2025-05-20

### Fixed

- AttributeError: module 'mlflow' has no attribute 'server' ([#134](https://github.com/mlflow/mlflow-go-backend/issues/134))

## [0.2.0] - 2025-01-27

### Added

- Support DeleteRegisteredModelAlias endpoint ([#48](https://github.com/mlflow/mlflow-go-backend/issues/48)).
- Support GetModelVersionByAlias endpoint ([#49](https://github.com/mlflow/mlflow-go-backend/issues/49)).

### Fixed

- Add alias for ui ([#112](https://github.com/mlflow/mlflow-go-backend/issues/112)).
- Add better message for file store uris ([#114](https://github.com/mlflow/mlflow-go-backend/issues/114)).
- Avoid SQL error log ([#119](https://github.com/mlflow/mlflow-go-backend/issues/119)).

## [0.1.0] - 2025-01-22

### Miscellaneous

- Initial release!
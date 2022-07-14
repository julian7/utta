# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

Changed:

* update dependencies
* switch to zap logger
* add log level filter as a global option
* remote: add circuit breaker to sleep for 30m (`--sleep`) after three (`--breaker`) consecutive
  connections being closed in 30s

## [v0.1.0] - Mar 28, 2021

Initial release

[Unreleased]: https://github.com/julian7/utta/
[v0.1.0]: https://github.com/julian7/utta/releases/tag/v0.1.0

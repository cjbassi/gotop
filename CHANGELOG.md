# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [3.0.0] - 2019-02-22

### Added

- Add vice colorscheme [#115]

### Changed

- Change `-v` cli option to `-V` for version
- Revert back to using the XDG spec on macOS

### Fixed

- Fix disk I/O [#114] [#116]

## [2.0.2] - 2019-02-16

### Fixed

- Fix processes on macOS not showing when there's a space in the command name [#107] [#109]

[#107]: https://github.com/cjbassi/gotop/issues/107
[#114]: https://github.com/cjbassi/gotop/issues/114

[#109]: https://github.com/cjbassi/gotop/pull/109
[#115]: https://github.com/cjbassi/gotop/pull/115
[#116]: https://github.com/cjbassi/gotop/pull/116

[Unreleased]: https://github.com/cjbassi/gotop/compare/3.0.0...HEAD
[3.0.0]: https://github.com/cjbassi/gotop/compare/2.0.2...3.0.0
[2.0.2]: https://github.com/cjbassi/gotop/compare/2.0.1...2.0.2

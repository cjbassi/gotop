# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> **Types of changes**:
>
> - **Added**: for new features.
> - **Changed**: for changes in existing functionality.
> - **Deprecated**: for soon-to-be removed features.
> - **Removed**: for now removed features.
> - **Fixed**: for any bug fixes.
> - **Security**: in case of vulnerabilities.

## [3.4.2] - ??

### Added

- Device data export via HTTP. If run with the `--export :2112` flag (`:2112`
  is a port), metrics are exposed as Prometheus metrics on that port.
- A battery gauge as a `power` widget; battery as a bar rather than
  a histogram.
- Temp widget displays degree symbol (merged from BartWillems, thanks
  also fleaz)
- Support for (device) plugins, and abstracting devices from widgets. This
  allows adding functionality without adding bulk.

### Fixed

- Keys not controlling process widget, #59
- The one-column bug, #62

## [3.3.2] - 2020-02-26

Bugfix release.

### Fixed

- #15, crash caused by battery widget when some accessories have batteries
- #57, colors with dashes in the name not found.
- Also, cjbassi/gotop#127 and cjbassi/gotop#130 were released back in v3.1.0.

## [3.3.1] - 2020-02-18

- Fixed: Fixes a layout bug where, if columns filled up early, widgets would be
  consumed but not displayed.
- Fixed: Rolled back dependency update on github.com/shirou/gopsutil; the new version
  has a bug that causes cores to not be seen.

## [3.3.0] - 2020-02-17

- Added: Logs are now rotated. Settings are currently hard-coded at 4 files of 5MB
  each, so logs shouldn't take up more than 20MB.  I'm going to see how many
  complain about wanting to configure these settings before I add code to do
  that.
- Added: Config file support. \$XDG_CONFIG_HOME/gotop/gotop.conf can now
  contain any field in Config.  Syntax is simply KEY=VALUE.  Values in config
  file are overridden by command-line arguments (although, there's a weakness
  in that there's no way to disable boolean fields enabled in the config).
- Changed: Colorscheme registration is changed to be less hard-coded.
  Colorschemes can now be created and added to the repo, without having to also
  add hard-coded references elsewhere.
- Changed: Minor code refactoring to support Config file changes has resulted
  in better isolation.

## [3.2.0] - 2020-02-14

Bug fixes & pull requests

- Fixed: Rowspan in a column loses widgets in later columns
- Fixed: Merged pull request for README clean-ups (theverything:add-missing-option-to-readme)
- Added: Merge Nord color scheme (jrswab:nordColorScheme)
- Added: Merge support for multiple (and filtering) network interfaces (mattLLVW:feature/network_interface_list)
- Added: Merge filtering subprocesses by substring (rephorm:filter)

## [3.1.0] - 2020-02-13

Re-homed the project after the original fork (trunk?) was marked as
unmaintained by cjbassi.

-  Changed: Merges @HowJMay spelling fixes
-  Added: Merges @markuspeloquin solarized themes
-  Added: Merges @jrswab additional kill terms
-  Added: Adds the ability to lay out the UI using a text file
-  Changed: the project filesystem layout to be more idiomatic

## [3.0.0] - 2019-02-22

### Added

- Add vice colorscheme [#115]

### Changed

- Change `-v` cli option to `-V` for version
- Revert back to using the XDG spec on macOS

### Fixed

- WIP fix disk I/O statistics [#114] [#116]

## [2.0.2] - 2019-02-16

### Fixed

- Fix processes on macOS not showing when there's a space in the command name [#107] [#109]

[#134]: https://github.com/cjbassi/gotop/issues/134
[#127]: https://github.com/cjbassi/gotop/issues/127
[#124]: https://github.com/cjbassi/gotop/issues/124
[#119]: https://github.com/cjbassi/gotop/issues/119
[#118]: https://github.com/cjbassi/gotop/issues/118
[#117]: https://github.com/cjbassi/gotop/issues/117
[#114]: https://github.com/cjbassi/gotop/issues/114
[#107]: https://github.com/cjbassi/gotop/issues/107
[#20]: https://github.com/cjbassi/gotop/issues/20

[#145]: https://github.com/cjbassi/gotop/pull/145
[#144]: https://github.com/cjbassi/gotop/pull/144
[#130]: https://github.com/cjbassi/gotop/pull/130
[#129]: https://github.com/cjbassi/gotop/pull/129
[#128]: https://github.com/cjbassi/gotop/pull/128
[#121]: https://github.com/cjbassi/gotop/pull/121
[#120]: https://github.com/cjbassi/gotop/pull/120
[#116]: https://github.com/cjbassi/gotop/pull/116
[#115]: https://github.com/cjbassi/gotop/pull/115
[#112]: https://github.com/cjbassi/gotop/pull/112
[#109]: https://github.com/cjbassi/gotop/pull/109

[Unreleased]: https://github.com/cjbassi/gotop/compare/3.0.0...HEAD
[3.0.0]: https://github.com/cjbassi/gotop/compare/2.0.2...3.0.0
[2.0.2]: https://github.com/cjbassi/gotop/compare/2.0.1...2.0.2

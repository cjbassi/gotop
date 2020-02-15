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

## [3.3.0] - 

- Changed: Logs are now rotated. Settings are currently hard-coded at 4 files of 5MB
  each, so logs shouldn't take up more than 20MB.  I'm going to see how many
  complain about wanting to configure these settings before I add code to do
  that.

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

## [Unreleased]

### Added

- Add snap package [#119] [#120] [#121]
- Process list scroll indicator [#127] [#130]
- Preliminary OpenBSD support [#112] [#117] [#118]

### Fixed

- Fix process localization issues on macOS [#124]
- Fix miscellaneous issues on FreeBSD [#134] [#145]
- Fix spelling of "Tx" to "TX" [#129]
- Rerender statusbar on every tick [#128]

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

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

## [4.0.0] PENDING

**Command line options have changed.**

### Added

- Adds support for system-wide configurations.  This improves support for package maintainers.
- Help function to print key bindings.
- Help prints locations of config files (color schemes & layouts).
- Help prints location of logs.
- CLI option to scale out (#84).
- Ability to report network traffic rates as mbps (#46).
- Ignore lines matching `/^#.*/` in layout files.
- Instructions for Gentoo (thanks @tormath1!)
- Graph labels that don't fit (vertically) in the window are now drawn in additional columns (#40)
- Adds ability to filter reported temperatures (#92)
- Command line option to list layouts, paths, colorschemes, hotkeys, and filterable devices
- Adds ability to write out a configuration file
- Adds a command for specifying the configuration file to use
- Added contribution from @wcdawn for building on machines w/ no Go/root access

### Changed

- Log files stored in \$XDG_CACHE_HOME; DATA, CONFIG, CACHE, and RUNTIME are the only directories specified by the FreeDesktop spec.
- Extensions are now built with a build tool; this is an interim solution until issues with the Go plugin API are resolved.
- Command line help text is cleaned up.
- Version bump of gopsutil
- Prometheus client replaced by [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics). This eliminated 6 indirect package dependencies and saved 3.5MB (25%) of the compiled binary size.
- Relicensed to MIT-3

### Removed

- configdir, logdir, and logfile options in the config file are no longer used.  gotop looks for a configuration file, layouts, and colorschemes in the following order: command-line; `pwd`; user-home, and finally a system-wide path.  The paths depend on the OS and whether XDG is in use.
- Removes the deprecated `--minimal` and `--battery` options.  Use `-l minimal` and `-l battery` instead.

### Fixed

- Help & statusbar don't obey theme (#47).
- Fix help text layout.
- Merged fix from @markuspeloquin for custom color scheme loading crash
- Memory line colors were inconsistently assigned (#91)
- The disk code was truncating values instead of rounding (#90)
- Temperatures on Darwin were all over the place, and wrong (#48)
- Config file loading from `~/.config/gotop` wasn't working
- There were a number of minor issues with the config file that have been cleaned up.

## [3.5.2] - 2020-04-28

### Fixed

- Fixes (an embarrasing) null map bug on FreeBSD (#94)

## [3.5.1] - 2020-04-09

This is a bug fix release.

### Fixed

- Removes verbose debugging unintentionally left in the code (#85)
- kitchensink referenced by, but not included in binary is now included (#72)
- Safety check prevents uninitialized colorscheme registry use
- Updates instructions on where to put colorschemes and layouts (#75)
- Trying to use a non-installed plugin should fail, not silently continue (#77)

### Changed

- Improved documentation about installing layouts and colorschemes

## [3.5.0] - 2020-03-06

The version jump from 3.3.x is due to some work in the build automation that necessitated a number of bumps to test the build/release, and testing compiling plugins from github repositories.

### Added

- Device data export via HTTP. If run with the `--export :2112` flag (`:2112`
  is a port), metrics are exposed as Prometheus metrics on that port.
- A battery gauge as a `power` widget; battery as a bar rather than
  a histogram.
- Temp widget displays degree symbol (merged from BartWillems, thanks
  also fleaz)
- Support for (device) plugins, and abstracting devices from widgets. This
  allows adding functionality without adding bulk. See the [plugins decision wiki section](https://github.com/xxxserxxx/gotop/wiki/Plugins-in-gotop) for more information.

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

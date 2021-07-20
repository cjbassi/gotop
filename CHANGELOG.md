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

## [4.1.2] 2021-07-20

### Added

- Several folks contributed to building on new Apple silicon (@clandmeter,
  @areese, and @nickcorin). This took a distressingly long time for me to merge;
  it required updating and testing a newer cross-compiling CGO, and I'm timid
  when it comes to releasing stuff I can't test.
- French and Russion translations (thank you @lourkeur and @talentlessguy!)
- nvidia support merged in from extension
- remote support merged in from extension
- Spanish translation (thanks to @donPatino & @lourkeur)
- There's a link to the github project in the help text now

### Changed

- Upgrade to Go 1.16. This eliminates go:generate for the language files, which
  means gotop no longer builds with Go < 1.16. It does make things easier for
  translators and merging.
- The [remote monitoring documentation](https://github.com/xxxserxxx/gotop/blob/master/docs/remote-monitoring.md) is a little better.

### Fixed

- Extra spaces in help text (#167)
- Crash with German translation (#166)
- Bad error message for missing layouts (#164)
- @JonathanReeve, @joinemm, and @plgruener contributed typo and mis-translation fixes
- The remote extension was ignoring config-file settings (no ticket #)

## [4.1.1] 2021-02-03 

### Added

- Show available translations in help text (#157)

### Changed

- Add more badges in README
- Replaces a dependency with a fork, because `go get` -- unlike `go build` -- ignores `replace` directives in `go.mod`
- Adds links to the extension projects in the README (#159)
- Missing thermal sensors on FreeBSD were being reported as errors, when they aren't. (#152)
- Small performance optimization
- github workflow changes to improve failure modes
- Bumped `gopsutils` to v3.20.12
- Bumped `battery` to v0.10.0
- Debug logging was left on (again) causing chatter in logs

### Fixed

- No temperatures on Raspberry Pi (#6)
- CPU name sorting in load widget (#161)
- The status bar got lost at some point; it's back
- Errors from any battery prevented display of all battery information


## [4.1.0] 2021-01-25

The minor version bump reflects the addition of I18N. If you are using one of the languages that has a translation, and your environment is set to that language, the UI will be different.  Translations are very welcome!

Thanks to the people who submitted PRs and translations to this release.

### Added

- Adds multilingual support.  German, Chinese (zh_CN), Esperanto (#120)

### Changed

- The uploaded license was a 2-clause BSD, which is functionally equivalent; however, since the contributor agreement was for MIT, to make it clean the uploaded license file was changed to the Festival variant of MIT. (#147)
- Per-process CPU use was averaged over the entire process lifetime.  While more of a semantic difference than a bug, it was a unintuitive and not particularly useful. CPU averages are now weighted moving averages over time, with more recent use having more weight.
- iSMC was still in the code; iSMC violates the MIT license, and this has been cleaned out.
- Versions are now embedded during the package build, rather than being hard-coded. More info in #140

### Fixed

- No temperatures displayed (#130), a recurring issue.
- Cannot show the CPU usages of the processes (#135)
- power widget consumes all RAM (#134)
- Disk usage not showing up at all (#27)
- Missing CPU: Lists 7 of 8 expected (#19)


## [4.0.1] 2020-06-08

**Darwin-only release**

### Changed

- The change to remove GPL dependencies did not remove *all* dependencies. This corrects that (#131)
 

## [4.0.0] 2020-06-07

**Command line options have changed.**

### Added

- Adds support for system-wide configurations.  This improves support for package maintainers.
- Help function to print key bindings, widgets, layouts, colorschemes, and paths
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
- Merged cmatsuoka's console font contribution
- Added contribution from @wcdawn for building on machines w/ no Go/root access

### Changed

- Log files stored in \$XDG_CACHE_HOME; DATA, CONFIG, CACHE, and RUNTIME are the only directories specified by the FreeDesktop spec.
- Extensions are now built with a build tool; this is an interim solution until issues with the Go plugin API are resolved.
- Command line help text is cleaned up.
- Version bump of gopsutil
- Prometheus client replaced by [VictoriaMetrics/metrics](https://github.com/VictoriaMetrics/metrics). This eliminated 6 indirect package dependencies and saved 3.5MB (25%) of the compiled binary size.
- Relicensed to MIT-3 (see [#36](https://github.com/xxxserxxx/gotop/issues/36))

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
- Compile errors on FreeBSD due to golang.org/x/sys API breakages
- Key bindings now work in FreeBSD (#95)
- Only report battery sensor errors once (reduce noise in the log, #117)
- Fixes a very small memory leak from the spark and histograph widgets (#128)

## [3.5.3] - 2020-05-30

The FreeBSD bugfix release. While there are non-FreeBSD fixes in here, the focus was getting gotop to work properly on FreeBSD.

### Fixed

- Address FreeBSD compile errors resulting to `golang.org/x/sys` API breakages
- Key bindings now work in FreeBSD (#95)
- Eliminate repeated logging about missing sensor data on FreeBSD VMs (#97)
- Investigated #14, a report about gotop's memory not matching `top`'s numbers, and came to the conclusions that (a) `gotop` is more correct in some cases (swap) than `top`, and (b) that the metric `gotop` is using (`hw.physmem`) is probably correct -- or that there's no obviously superior metric. So no change.

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

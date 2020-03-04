---
name: Bug report
about: Template to report bugs.
---

<!-- Please search existing issues to avoid creating duplicates. -->
<!-- Also please test using the latest build to make sure your issue has not already been fixed. -->

##### gotop version:
`gotop -V`, or if built from source, `git rev-parse HEAD`
##### OS/Arch:
Linux: `uname -or`, OSX: `sw_vers`; Windows: `systeminfo | findstr /B /C:"OS Name" /C:"OS Version"`
##### Terminal emulator: 
e.g. iTerm, kitty, xterm, PowerShell
##### Any relevant hardware info:
If the issue is clearly related to a specific piece of hardware, e.g., the network
##### tmux version:
`tmux -V`, if using tmux

Also please copy or attach `~/.local/state/gotop/errors.log` if it exists and contains logs:

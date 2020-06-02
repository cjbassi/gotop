# Config file

Most command-line settings can be persisted into a configuration file. The config file is named `gotop.conf` and can be located in several places. The first place gotop will look is in the current directory; after this, the locations depend on the OS and distribution. On Linux using XDG, for instance, the home location of `~/.config/gotop/gotop.conf` is the second location. The last location is a system-wide global location, such as `/etc/gotop/gotop.conf`. The `-h` help command will print out all of the locations, in order. Command-line options override values in any config files, and only the first config file found is loaded.

A configuration file can be created using the `--write-config` command-line argument. This will try to place the config file in the home config directory (the second location), but if it's unable to do so it'll write a file to the current directory.

Config file changes can be made by combining command-line arguments with `--write-config`. For example, to persist the `solarized` theme, call:

```
gotop -c solarized --write-config
```


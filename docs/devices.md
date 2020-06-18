# Device filtering

Some devices have quite a number of data points; on OSX, for instance, there are dozens of temperature readings. These can be filtered through a configuration file.  There is no command-line argument for this filter.

The list will grow, but for now the only device that supports filtering is the temperature widget.  The configuration entry is called `temperature`, and it contains an exact-match list of comma-separated values with no spaces.  To see the list of valid values, run gotop with the `--list devices` command.  Gotop will print out the type of device and the legal values.  For example, on Linux:

```
$ gotop --list devices
Temperatures:
        acpitz
        nvme_composite
        nvme_sensor1
        nvme_sensor2
        pch_cannonlake
        coretemp_packageid0
        coretemp_core0
        coretemp_core1
        coretemp_core2
        coretemp_core3
        ath10k_hwmon
```
You might then add the following line to the config file.  First, find where gotop looks for config files:
```
➜  gotop --list paths
Loadable colorschemes & layouts, and the config file, are searched for, in order:
/home/ser/workspace/gotop.d/gotop
/home/ser/.config/gotop
/etc/xdg/gotop

The log file is in /home/ser/.cache/gotop/errors.log
```
So you might use `${HOME}/.config/gotop/gotop.conf`, and add (or modify) this line:
```
temperatures=acpitz,coretemp_core0,ath10k_hwmon
```
This will cause the temp widget to show only four of the eleven temps.


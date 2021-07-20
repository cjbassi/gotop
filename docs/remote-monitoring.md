# Remote monitoring extension for gotop


Show data from gotop running on remote servers in a locally-running gotop. This allows gotop to be used as a simple terminal dashboard for remote servers.

![Screenshot](/assets/screenshots/fourby.png)


## Configuration

gotop exports metrics on a local port with the `--export <port>` argument. This is a simple, read-only interface with the expectation that it will be run behind some proxy that provides security.  A gotop built with this extension can read this data and render it as if the devices being monitored were on the local machine.

On the local side, gotop gets the remote information from a config file; if all you have is a single remote machine to monitor, the parameters can be passed on the command line. For more than one remote, a config file is needed. The recommended approach is to create a remote-specific config file, and then run gotop with the `-C <remote-config-filename>` option. The plan is to add two functions that will enhance this feature: disabling the UI on the remote machine (allowing gotop to be forked into the background); and disabling local metrics to focus a gotop instance on remote machines. Also planned are a data transfer optimization and increasing the metrics that can be monitored.

Two options are available for each remote server; one of these, the connection URL, is required.  The format of the configuration keys are: `remote-SERVERNAME-url` and `remote-SERVERNAME-refresh`; `SERVERNAME` can be anything -- it doesn't have to reflect any real attribute of the server, but it will be used in widget labels for data from that server.  For example, CPU data from `remote-Jerry-url` will show up as `Jerry-CPU0`, `Jerry-CPU1`, and so on; memory data will be labeled `Jerry-Main` and `Jerry-Swap`.  If the refresh rate option is omitted, it defaults to 1 second.


### An example

One way to set this up is to run gotop behind [Caddy](https://caddyserver.com). The `Caddyfile` would have something like this in it:

```
gotop.myserver.net {
        basicauth / gotopusername supersecretpassword
        reverse-proxy / http://localhost:8089
}
```

Caddy would then be responsible for authentication and encrypting the traffic.  Then, on the same machine run gotop in a persistent terminal session such as [tmux](https://github.com/tmux/tmux) with the following command:

```
gotop -x :8089
```

On a local machine, create a config file named `myserver.conf` with the following lines:

```
remote-myserver-url=https://gotopusername:supersecretpassword@gotop.myserver.net/metrics
remote-myserver-refresh=2
```

Note the `/metrics` at the end -- don't omit that, and don't strip it in Caddy.  The refresh value is in seconds. Run gotop with:

```
gotop -C myserver.conf
```

and you should see your remote server sensors as if it were running on your local machine.

You can add as many remote servers as you like in the config file; just follow the naming pattern.

## Why

This can combine multiple servers into one view, which makes it more practical to use a terminal-based monitor when you have more than a couple of servers, or when you don't want to dedicate an entire wide-screen monitor to a bunch of gotop instances. It's simple to set up, configure, and run, and reasonably resource efficient.

## How

Since v3.5.2, gotop's been able to export its sensor data as [Prometheus](https://prometheus.io/) metrics using the `--export` flag.  Prometheus has the advantages of being simple to integrate into clients, and a nice on-demand design that depends on the *aggregator* pulling data from monitors, rather than the clients pushing data to a server. In essence, it inverts the client/server relationship for monitoring/aggregating servers and the things it's monitoring. In gotop's case, it means you can turn on `-x` and not have it impact your gotop instance at all, until you actively poll it.  It puts the control on measurement frequency in a single place -- your local gotop. It means you can simply stop your local gotop instance (e.g., when you go to bed) and the demand on the servers you were monitoring drops to 0. 

On the client (local) side, sensors are abstracted as devices that are read by widgets, and we've simply implemented virtual devices that poll data from remote Prometheus instances. At a finer grain, there's a single process spawned for each remote server that periodically polls that server and collects the information.  When the widget updates and asks the virtual device for data, the device consults the cached data and provides it as the measurement.

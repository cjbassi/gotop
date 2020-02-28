% Plugins


# Extensions

- Plugins will supply an `Init()` function that will call the appropriate
  `Register\*()` functions in the `github.com/xxxserxxx/gotop/devices` package.
- `devices` will supply:
    - RegisterCPU (opt)
        - Counts (req)
        - Percents (req)
    - RegisterMem (opt)
    - RegisterTemp (opt)
    - RegisterShutdown (opt)

# gotop

- Command line -P, comma separated list of plugin .so
- gotop will look in `pwd` and then in \$XDG_CONFIG_HOME/gotop
- When loaded, gotop will call lib#Init()

When exited cleanly, gotop will call all registered shutdown functions.

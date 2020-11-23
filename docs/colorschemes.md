# Colorschemes

gotop ships with a few colorschemes which can be set with the `-c` flag followed by the name of one. You can find all the colorschemes in the [colorschemes folder](../colorschemes).

To make a custom colorscheme, check out the [template](../colorschemes/template.go) for instructions and then use [default.json](../colorschemes/default.json) as a starter. Then put the file at `~/.config/gotop/<name>.json` and load it with `gotop -c <name>`. Colorschemes PR's are welcome!

To list all built-in color schemes, call:

```
gotop --list colorschemes
```


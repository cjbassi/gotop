# gotop NVidia extension

Provides NVidia GPU data to gotop

To enable it, either run gotop with the `--nvidia` flag, or add the line `nvidia=true` to `gotop.config`.

## Dependencies

- [nvidia-smi](https://wiki.archlinux.org/index.php/NVIDIA/Tips_and_tricks#nvidia-smi)

## Configuration

The refresh rate of NVidia data is controlled by the `nvidia-refresh` parameter in the configuration file.  This is a Go `time.Duration` format, for example `2s`, `500ms`, `1m`, etc.

## Alternatives to test

https://github.com/mindprince/gonvml

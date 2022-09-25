# Overview
`rocksdbctl` is a command line program to access rocksdb data.

# Installing

```sh
git clone https://github.com/dxhome/rocksdbctl
cd rocksdbctl
go build -o rocksdbctl
```

# Usage
## Dump Keys

```sh
# ./rocksdbctl dump -h
dump all key/value

Usage:
  rocksdbctl dump [flags]

Flags:
  -h, --help            help for dump
  -l, --limit int       dump limited keys (default 50)
      --prefix string   key prefix

Global Flags:
  -d, --debug         debug mode
  -p, --path string   rocksdb path
```

## Get Key

```sh
# ./rocksdbctl get -h
get a key/value

Usage:
  rocksdbctl get <key> [flags]

Flags:
  -b, --byte   disply mode in byte, default display mode is string
  -h, --help   help for get

Global Flags:
  -d, --debug         debug mode
  -p, --path string   rocksdb path
```

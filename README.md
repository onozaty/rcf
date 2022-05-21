# rcf

[![GitHub license](https://img.shields.io/github/license/onozaty/rcf)](https://github.com/onozaty/rcf/blob/main/LICENSE)
[![Test](https://github.com/onozaty/rcf/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/rcf/actions/workflows/test.yaml)

`rcf` is a CLI tool for replacing file contents.

## Usage

```
$ rcf -i input -s before -t after -o output
```

The arguments are as follows.

```
Usage: rcf -i INPUT [-r REGEX | -s STRING] -t REPLACEMENT [--escape] [--recursive] [-c CHARSET] [-o OUTPUT | --overwrite]

Flags
  -i, --input string         Input file/dir path.
  -r, --regex string         Target regex.
  -s, --string string        Target string.
  -t, --replacement string   Replacement.
  -e, --escape               Enable escape sequence.
  -R, --recursive            Recursively traverse the input dir.
  -c, --charset string       Charset. (default "UTF-8")
  -o, --output string        Output file/dir path.
  -O, --overwrite            Overwrite the input file.
  -h, --help                 Help.
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/rcf/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)


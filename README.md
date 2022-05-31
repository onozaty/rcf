# rcf

[![GitHub license](https://img.shields.io/github/license/onozaty/rcf)](https://github.com/onozaty/rcf/blob/main/LICENSE)
[![Test](https://github.com/onozaty/rcf/actions/workflows/test.yaml/badge.svg)](https://github.com/onozaty/rcf/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/onozaty/rcf/branch/main/graph/badge.svg?token=MAL9FJ1QW3)](https://codecov.io/gh/onozaty/rcf)

`rcf` is a CLI tool for replacing file contents.

## Usage

```
$ rcf -i input.txt -s before -t after -o output.txt
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

### Specify replacement

The target can be a regular expression or a string.

The regular expression is specified by `-r`.

```
$ rcf -i input.txt -r "[0-9]" -t "" -o output.txt
```

You can also use the capture group as `-t`.  
The following is a method for specifying "N" at the beginning of a number.

```
$ rcf -i input.txt -r "([0-9]+)" -t "N$1" -o output.txt
```

Please refer to the following for the syntax of regular expressions.

* https://pkg.go.dev/regexp/syntax

String is specified by `-s`.

```
$ rcf -i input.txt -s a -t z -o output.txt
```

To treat backslash as an escape sequence, specify `-e`.

```
$ rcf -i input.txt -s "\u3000" -e -t "" -o output.txt
```

### Input / Output

If specified with `-i`, only the specified file will be processed.

```
$ rcf -i input.txt -s a -t z -o output.txt
```

If a directory is specified, files under the directory are processed.  

```
$ rcf -i in_dir -s a -t z -o out_dir
```

The default is to process only files directly under the specified directory.  
If `-r` is specified, subdirectories are processed recursively.

```
$ rcf -i in_dir -r -s a -t z -o out_dir
```

Use `-o` to specify the output destination.  
To rewrite the input file itself, use `-O`.

```
$ rcf -i input.txt -s before -t after -O
```

### Charset

When processing non UTF-8 files, specify the Charset with `-c`.

```
$ rcf -i input.txt -s a -t z -o output.txt -c sjis
```

Charset must be one that can be specified in `htmlindex.Get`.

* https://pkg.go.dev/golang.org/x/text/encoding/htmlindex#Get

A special charset is `binary`.  
If `binary` is specified, it can be treated as a hexadecimal character.  
A hexadecimal character represents a byte with three characters prefixed by `x`, such as `x00` or `xFF`.

To remove two consecutive bytes, such as 0x00 0x01, specify as follows

```
$ rcf -i input.txt -s x00x01 -t "" -o output.txt -c binary
```

## Install

You can download the binary from the following.

* https://github.com/onozaty/rcf/releases/latest

## License

MIT

## Author

[onozaty](https://github.com/onozaty)


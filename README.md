# Text Builder v0.2

Text Builder is a small utility which combines local and remote files into single text file.

## Status

This tool is in early stage development but stable enough to run larger benchmark sets.
missing features will be added as needed, pull requests are welcome ;)

## Installing

You need go 1.0+ (1.5+ is recommended)

If you just want to run this tool:

```
go get github.com/janlay/text-builder
```

Build from source code:

```
git clone git://github.com/janlay/text-builder.git
cd text-builder
go build
```

Run `text-builder -help` to see usage.

## Usage

Basic usage is quite simple:
```
text-builder
```
which uses STDIN as input and outputs to STDOUT.

For most cases, you may want to:
```
text-builder -index /path/to/index -output /path/to/output
```

Build the file and skip some #include directives:
```
text-builder -index /path/to/index -output /path/to/output -skip=a.txt,b.txt
```

## Layout of Index and Other files
See `./examples` for more details.
Check this [gist](https://gist.github.com/janlay/b57476c72a93b7e622a6) to understand how Text Builder builds file from several sources.

## Useful Tips

* Use `#output /path/to/file` in the first line of the index file to override `-output` argument.
* Use `#include /path/to/file` to include a local file.
* Use `#include http://host/path` to include a remote file.
* You may want to use relative file paths for `#output` and `#include`. These paths are related to current file, not working directory. This manner works like relative paths in .css files.
* `text-builder -index /path/to/index > output.txt` is equal to `text-builder -index /path/to/index -output output.txt`

## License

This software is licensed under the MIT License.

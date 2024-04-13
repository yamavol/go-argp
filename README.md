
# go-argp

Command line argument parser library, using similar syntax to the GNU argp interface. 

    command -abCdE --limit 10 --level=123 arg0 arg1

## Features

- Supports GNU style option ([syntax](./syntax.md))
- Testable API
- No conversion, only string parsing
- Built in method for generating nice formatted help message.
- Easy integration: only ~2 files

## Install

    go get github.com/yamavol/go-argp

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/yamavol/go-argp"
)

var options = []argp.Option{
	{Doc: "OPTIONS:"},
	{Short: 'a', Long: "", Doc: "doc"},
	{Short: 'b', Long: "bb", Doc: "doc"},
	{Short: 's', Long: "silent", Doc: "doc"},
	{Short: 'q', Long: "", Flags: argp.OPTION_ALIAS},
	{Short: 'o', Long: "output", ArgName: "<file>", Doc: "doc"},
	{Short: '1', Long: "", Doc: "doc"},
	{Doc: ""},
	{Short: 'h', Long: "help", Doc: "print help"},
	{Short: 'V', Long: "version", Doc: "print version"},
}

func main() {
	// argp.Parse(options) == argp.ParseArgs(options, os.Args[1:])
	result, err := argp.Parse(options)

	if err != nil {
		fmt.Printf("error: %s", err)
		return
	}

	if result.HasOpt("help") {
		// Prints usage and option list.
        // ```
		//   Usage: cmd [options...] ARG1 ARG2
		//   OPTIONS:
		//    -a                        doc
		//    -b, --bb                  doc
		// ```
		// argp.PrintOptList() prints without the usage line
		argp.PrintUsage(os.Stdout, options, "cmd", "ARG1 ARG2...")
		return
	}
	if result.HasOpt("version") {
		fmt.Println("0.0.0")
		return
	}
	// result.Options is a list of [Option]s found in order
	for _, opt := range result.Options {
		fmt.Printf("%s: %s\n", opt.InputString, opt.Optarg)
	}
	// result.Args is a list of non-option Arguments
	fmt.Print(result.Args)
}

```

## Guidance

The golang's standard way to parse options is to use the `flag` package. Actually it is good enough to get the job done.

The syntax of `flag` package differs from other well known tools. Most UNIX and GNU tools distinguishes short option and long option. Each have different syntaxes; for example, = sign is only usable in long options. `flag` mixes these rules, but it might be better to follow the major convention. 

Other packages offers better usability and readability. Such features include subcommands, multiple arguments, string conversion, string splitting. 

## License

MIT


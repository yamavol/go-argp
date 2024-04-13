package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	// argp.Parse(options) or argp.ParseArgs(options, os.Args[1:])
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
		argp.PrintUsage(os.Stdout, options, filepath.Base(os.Args[0]), "ARG1 ARG2...")
		return
	}
	if result.HasOpt("version") {
		fmt.Println("1.0.0")
		return
	}
	// result.Options is a list of [Option]s found in order
	for _, opt := range result.Options {
		fmt.Printf("%s: %s\n", opt.InputString, opt.Optarg)
	}
	// result.Args is a list of non-option Arguments
	fmt.Print(result.Args)
}

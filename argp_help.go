package argp

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

func empty_rune(c rune) bool {
	return c == 0 || c == ' '
}

func empty_str(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}

// Prints the help message to the [io.Writer]
func PrintUsage(w io.Writer, options []Option, cmd string, arg string) {
	fmt.Fprintf(w, "Usage: %s [options...] %s\n", cmd, arg)
	PrintOptList(w, options)
}

// Prints the help message to the [io.Writer]. This help message only contains
// the option list.
func PrintOptList(w io.Writer, options []Option) {
	var pOptReal *Option
	var pOptAliases []*Option
	for lc := 0; lc < len(options); lc++ {

		opt := options[lc]

		// list the following alias commands, increment the loop counter
		if opt.Flags&OPTION_ALIAS == 0 {
			pOptReal = &options[lc]

			for off := 1; lc+off < len(options); off++ {
				if options[lc+off].Flags&OPTION_ALIAS != 0 {
					pOptAliases = append(pOptAliases, &options[lc+off])
				} else {
					lc += (off - 1)
					break
				}
			}
		} else {
			continue
		}

		if opt.Flags&OPTION_HIDDEN > 0 {
			// don't print hidden option
			continue
		} else if empty_rune(opt.Short) && empty_str(opt.Long) && len(opt.Doc) > 0 {
			// print category header
			fmt.Fprintln(w, opt.Doc)
		} else {
			// print options and its decriptions
			left := sprintfOptions(pOptReal, pOptAliases)

			docRows := strings.Split(pOptReal.Doc, "\n")

			for _, row := range docRows {
				fmt.Fprintf(w, "%-25s  %s\n", left, row)
				left = ""
			}

			pOptAliases = nil
		}
	}
}

type argFmt int

const (
	argFmtNone argFmt = iota
	argFmtDefault
	argFmtLongOptional
	argFmtShortOptional
)

// Return a formatted string for printing ArgName in the help message.
// Such examples are " ARG" or "[=ARG]" or "[ARG"].
// the output is determined by the argfmt argument.
func sprintfArg(arg string, argfmt argFmt) string {
	if empty_str(arg) {
		return ""
	}
	switch argfmt {
	case argFmtNone:
		return ""
	case argFmtDefault:
		return fmt.Sprintf(" %s", arg)
	case argFmtLongOptional:
		return fmt.Sprintf("[=%s]", arg)
	case argFmtShortOptional:
		return fmt.Sprintf("[%s]", arg)
	default:
		return "invalid argfmt"
	}
}

// Return a single dashed option string with an argument attached
func sprintfShort(short string, arg string, argfmt argFmt) string {
	return fmt.Sprintf("-%s%s", short, sprintfArg(arg, argfmt))
}

// Return a double dashed option string with an argument attached
func sprintfLong(long string, arg string, argfmt argFmt) string {
	return fmt.Sprintf("--%s%s", long, sprintfArg(arg, argfmt))
}

// Returns a formatted string of options with argument names attached:
//
//	"-o, --output <file>"
//
// The output format depends on how short/long options
// provided, and their attributes (FLAGS). It tries to mimic the original
// GNU argp output
func sprintfOptions(optReal *Option, optAlias []*Option) string {
	var runes []string
	var longs []string
	var optList = append([]*Option{optReal}, optAlias...)
	for _, opt := range optList {
		if !empty_rune(opt.Short) {
			runes = append(runes, string(opt.Short))
		}
		if !empty_str(opt.Long) {
			longs = append(longs, string(opt.Long))
		}
	}

	var buf bytes.Buffer

	// indent
	buf.WriteString(" ")

	// print short options.
	if len(runes) > 0 {
		list := []string{}
		for _, c := range runes {
			token := ""
			if len(longs) > 0 || empty_str(optReal.ArgName) {
				// if long option name is defined, skip short option arguments
				token = sprintfShort(c, "", argFmtNone)
			} else {
				// print short name and argName
				argfmt := argFmtDefault
				if optReal.Flags&OPTION_ARG_OPTIONAL > 0 {
					argfmt = argFmtShortOptional
				}
				token = sprintfShort(c, optReal.ArgName, argfmt)
			}
			list = append(list, token)
		}
		buf.WriteString(strings.Join(list, ", "))
	} else {
		// indent if no short option
		buf.WriteString("    ")
	}

	if len(runes) > 0 && len(longs) > 0 {
		// print a separator between the short and long option
		buf.WriteString(", ")
	}

	// print long options
	if len(longs) > 0 {
		list := []string{}
		for _, long := range longs {
			token := ""
			if empty_str(optReal.ArgName) {
				token = sprintfLong(long, "", argFmtNone)
			} else {
				// print long name and argName
				argfmt := argFmtDefault
				if optReal.Flags&OPTION_ARG_OPTIONAL > 0 {
					argfmt = argFmtLongOptional
				}
				token = sprintfLong(long, optReal.ArgName, argfmt)
			}
			list = append(list, token)
		}
		buf.WriteString(strings.Join(list, ", "))
	}

	return buf.String()
}

// package argp provides functions to parse command line options
package argp

import (
	"fmt"
	"os"
	"strings"
)

const (
	// Mark this option's argument as optional.
	// If the argument is optional, the argument must be provided in attached
	// style, otherwise it will raise an error. E.g. -o<ARG> or --option=<ARG>
	OPTION_ARG_OPTIONAL = 0x1

	// Hide this option from the help message
	OPTION_HIDDEN = 0x2

	// Mark this option as an alias of the previous non-alias option.
	// An alias option will be resolved to a non-alias option.
	OPTION_ALIAS = 0x4

	// [Private] Mark this option as "Non option". This flag is used to
	// return the non-option argument to support option reordering
	_OPTION_NON_OPTION_ARG = 0x20

	ErrInvalid = "invalid option"
	ErrMissing = "option requires an argument"
	ErrTooMany = "option takes no arguments"
)

// Option struct represents a single option.
// An Option table or an array of Option is used to parse the string array,
// and to generate the help message.
type Option struct {
	Short   rune   // The short option name. Use alphanum, otherwise 0/SPACE.
	Long    string // The long option name. Set empty string if unused.
	ArgName string // The name of argument if it takes one.
	Flags   int    // Option flags
	Doc     string // Description, or a single line text for header/line
}

// Returns true if the short name or long name equals the argument
func (o *Option) Is(name string) bool {
	return !empty_str(name) && (name == string(o.Short) || name == o.Long)
}

// Error object implements error interface, and extends the option entry
// which raised an error.
type Error struct {
	Option
	Message string
}

func (e Error) Error() string {
	if !empty_str(e.Long) && !empty_rune(e.Short) {
		return fmt.Sprintf("%s: --%s (-%c)", e.Message, e.Long, e.Short)
	} else if !empty_str(e.Long) {
		return fmt.Sprintf("%s: --%s", e.Message, e.Long)
	} else {
		return fmt.Sprintf("%s: -%c", e.Message, e.Short)
	}
}

// Result is an individual successfully parsed option. It embeds the original
// option and the argument.
type Result struct {
	Option
	InputString string // The original string supplied in the argument
	Optarg      string // option argument
}

// Return Optarg with default string
func (p *Result) WithDefault(arg string) string {
	if p == nil {
		return arg
	} else if empty_str(p.Optarg) {
		return arg
	} else {
		return p.Optarg
	}
}

type ParseResult struct {
	Options []Result
	Args    []string
}

// Check if option with given name was specified
func (p *ParseResult) HasOpt(long string) bool {
	return len(p.GetOpts(long)) > 0
}

// Get the first option with given name. Returns nil if not found.
func (p *ParseResult) GetOpt(name string) *Result {
	opts := p.GetOpts(name)
	if len(opts) > 0 {
		return opts[0]
	} else {
		return nil
	}
}

// Get all options with given name, both long and short
func (p *ParseResult) GetOpts(name string) []*Result {
	var results []*Result
	for i, opt := range p.Options {
		if opt.Short == rune(name[0]) || opt.Long == name {
			results = append(results, &p.Options[i])
		}
	}
	return results
}

// Parse string array
func ParseArgs(options []Option, args []string) (ParseResult, error) {
	parser := parser{options: options, args: args}
	var result ParseResult
	for {
		opt, err := parser.next()
		if err != nil || opt == nil {
			result.Args = append(result.Args, parser.rest()...)
			return result, err
		}
		if opt.Flags&_OPTION_NON_OPTION_ARG > 0 {
			result.Args = append(result.Args, opt.Optarg)
		} else {
			result.Options = append(result.Options, *opt)
		}
	}
}

// Parse [os.Args] provided
func Parse(options []Option) (ParseResult, error) {
	return ParseArgs(options, os.Args[1:])
}

// parser extracts options one-by-one from the string array.
type parser struct {
	options []Option // user-defined option table (readonly)
	args    []string // user-provided argument list (readonly)
	optidx  int      // parse index
	subopt  int      // sub-index to parse short options
}

// extracts one short option from the arg array
func (p *parser) short() (*Result, error) {
	runes := []rune(p.args[p.optidx])

	c := runes[p.subopt]
	option := findShort(p.options, c)

	if option == nil {
		return nil, Error{Option{Short: c}, ErrInvalid}
	}

	cstr := string(c)

	if len(option.ArgName) == 0 {
		p.subopt++
		if p.subopt >= len(runes) {
			p.subopt = 0
			p.optidx++
		}
		return &Result{*option, cstr, ""}, nil
	}
	if option.Flags&OPTION_ARG_OPTIONAL == 0 {
		optarg := string(runes[p.subopt+1:])
		p.subopt = 0
		p.optidx++
		if optarg == "" {
			if p.optidx == len(p.args) {
				return nil, Error{*option, ErrMissing}
			}
			optarg = p.args[p.optidx]
			p.optidx++
		}
		return &Result{*option, cstr, optarg}, nil
	} else {
		optarg := string(runes[p.subopt+1:])
		p.subopt = 0
		p.optidx++
		return &Result{*option, cstr, optarg}, nil
	}
}

// extracts one short option from the arg array
func (p *parser) long() (*Result, error) {
	long := p.args[p.optidx][2:]

	eq := strings.IndexByte(long, '=')
	var optarg string
	var attached bool
	if eq != -1 {
		optarg = long[eq+1:]
		long = long[:eq]
		attached = true
	}

	option := findLong(p.options, long)
	if option == nil {
		return nil, Error{Option{Long: long}, ErrInvalid}
	}

	// consume one token here, after valid option was found
	p.optidx++

	if len(option.ArgName) == 0 { // No argument
		if attached {
			return nil, Error{*option, ErrTooMany}
		}
		return &Result{*option, long, ""}, nil
	}

	if option.Flags&OPTION_ARG_OPTIONAL == 0 {
		if !attached {
			if p.optidx >= len(p.args) {
				return nil, Error{*option, ErrMissing}
			}
			optarg = p.args[p.optidx]
			p.optidx++
		}
		return &Result{*option, long, optarg}, nil
	} else {
		return &Result{*option, long, optarg}, nil
	}
}

// extracts one option from the arg array
func (p *parser) next() (*Result, error) {
	if p.optidx >= len(p.args) {
		return nil, nil
	}

	arg := p.args[p.optidx]

	if p.subopt > 0 {
		return p.short() // continue parsing short options
	}

	if len(arg) < 2 || arg[0] != '-' {
		p.optidx++
		return makeArg(arg), nil
	}

	if arg == "--" {
		p.optidx++
		return nil, nil
	}

	if arg[:2] == "--" {
		return p.long()
	}
	if arg[:1] == "-" {
		p.subopt = 1
		return p.short()
	}
	p.optidx++
	return makeArg(arg), nil
}

func (p *parser) rest() []string {
	return p.args[p.optidx:]
}

func findLong(options []Option, long string) *Option {
	var oOptReal *Option
	for i, option := range options {
		if option.Flags&OPTION_ALIAS == 0 {
			oOptReal = &options[i]
		}
		if oOptReal != nil && option.Long == long {
			return oOptReal
		}
	}
	return nil
}

func findShort(options []Option, short rune) *Option {
	var pOptReal *Option
	for i, option := range options {
		if option.Flags&OPTION_ALIAS == 0 {
			pOptReal = &options[i]
		}
		if pOptReal != nil && option.Short == short {
			return pOptReal
		}
	}
	return nil
}

func makeArg(text string) *Result {
	return &Result{
		Option: Option{
			Flags: _OPTION_NON_OPTION_ARG,
		},
		Optarg: text,
	}
}

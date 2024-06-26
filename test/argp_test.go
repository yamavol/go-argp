package argp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yamavol/go-argp"
	"github.com/yamavol/go-argp/test/harness"
)

var options = []argp.Option{
	{Doc: "CATEGORY 000:"},
	{Short: 'a', Long: "aaa", Doc: "enable option a"},
	{Short: 'b', Long: "bbb", Doc: "enable option b"},
	{Short: 'c', Long: "ccc", Doc: "enable option c"},
	{Short: 'd', Long: "ddd", ArgName: "<ARG>"},
	{Short: 'e', Long: "eee", ArgName: "<ARG>", Flags: argp.OPTION_ARG_OPTIONAL},

	{Short: 'p', Long: "", Doc: "enable option p"},
	{Short: 'q', Long: "", Doc: "enable option q"},
	{Short: 'r', Long: "", Doc: "enable option r"},
	{Short: '1', Long: "", Doc: "enable option 1"},
	{Short: ' ', Long: "secret", Flags: argp.OPTION_HIDDEN, Doc: "hidden option"},

	{Doc: "CATEGORY 111:"},
	{Short: 'o', Long: "output", ArgName: "<buf>", Flags: 0, Doc: "specify output buffer"},
	{Short: 'x', Long: "xxxx", ArgName: "<arg>", Flags: 0, Doc: "enable option x"},
	{Short: 'f', Long: "file", ArgName: "<file>", Flags: 0, Doc: "file to open"},
	{Short: ' ', Long: "ffff", Flags: argp.OPTION_ALIAS},
	{Short: ' ', Long: "fgfg", Flags: argp.OPTION_ALIAS},
	{Short: 'F', Long: "    ", Flags: argp.OPTION_ALIAS},
	{Short: 'K', Long: "kind", ArgName: "<kind>", Flags: argp.OPTION_ARG_OPTIONAL, Doc: "specify kind"},

	{Doc: "CATEGORY 222:"},
	{Short: 'h', Long: "help", Doc: "print help and exit"},
	{Short: 'V', Long: "version", Doc: "print version and exit"},
}

func split(str string) []string {
	// TODO: split string with "" properly
	return strings.Split(str, " ")
}

func Test_Parse(t *testing.T) {
	args := split("-abc arg0 arg1 arg2")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "no error")
	harness.IsEqual(t, len(result.Options), 3, "result is not empty")
	harness.IsEqual(t, len(result.Args), 3, "rest is not empty")
}

func Test_ParseLong(t *testing.T) {
	args := split("--help --version -abc --kind -pqr -")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 9, "")
	harness.IsEqual(t, len(result.Args), 1, "")
}

func Test_OptionWithArgument(t *testing.T) {
	args := split("-x ARG")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "ARG", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	args = split("-xARG")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "ARG", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	args = split("--xxxx ARG")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "ARG", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	args = split("--xxxx=ARG")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "ARG", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	// option without the required argument
	args = split("-x")
	_, err = argp.ParseArgs(options, args)
	harness.IsNotNil(t, err, "")

	// long option without the required argument
	args = split("--xxxx")
	_, err = argp.ParseArgs(options, args)
	harness.IsNotNil(t, err, "")

	// long option without the required argument
	args = split("--xxxx -- -a -- -123")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 2, "")
	harness.IsEqual(t, len(result.Args), 1, "")
	harness.IsTrue(t, result.Options[0].Is("xxxx"), "")
	harness.IsEqual(t, result.Options[0].Optarg, "--", "")
	harness.IsTrue(t, result.Options[1].Is("a"), "")
	harness.IsEqual(t, result.Options[1].Optarg, "", "")
	harness.IsEqual(t, result.Args[0], "-123", "")
}

func Test_Version(t *testing.T) {
	args := split("--version")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, len(result.Args), 0, "")
}

func Test_DoubleDashTerminator(t *testing.T) {
	args := split("-p -q -r -- -s -t -u")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 3, "")
	harness.IsEqual(t, len(result.Args), 3, "")
}

func Test_SingleDashPlaceholder(t *testing.T) {
	args := split("-o-")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "-", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	args = split("-o -")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "-", "")
	harness.IsEqual(t, len(result.Args), 0, "")

	args = split("--output=-")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 1, "")
	harness.IsEqual(t, result.Options[0].Optarg, "-", "")
	harness.IsEqual(t, len(result.Args), 0, "")
}

func Test_Alias(t *testing.T) {
	args := split("-f input1.txt --file input2.txt --ffff input3.txt --fgfg=input4.txt")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsEqual(t, len(result.Options), 4, "")
	harness.IsEqual(t, result.Options[0].Long, "file", "")
	harness.IsEqual(t, result.Options[1].Long, "file", "")
	harness.IsEqual(t, result.Options[2].Long, "file", "")
	harness.IsEqual(t, result.Options[3].Long, "file", "")
	harness.IsEqual(t, result.Options[0].Optarg, "input1.txt", "")
	harness.IsEqual(t, result.Options[1].Optarg, "input2.txt", "")
	harness.IsEqual(t, result.Options[2].Optarg, "input3.txt", "")
	harness.IsEqual(t, result.Options[3].Optarg, "input4.txt", "")
	harness.IsEqual(t, len(result.Args), 0, "")
}

func Test_OptionsAPI(t *testing.T) {
	opts := argp.Option{
		Short: 'a',
		Long:  "aaa",
	}
	harness.IsTrue(t, opts.Is("a"), "")
	harness.IsTrue(t, opts.Is("aaa"), "")
	harness.IsFalse(t, opts.Is(""), "")
}

func Test_ResultAPI(t *testing.T) {
	args := split("--ddd ARG")
	result, err := argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsTrue(t, result.HasOpt("ddd"), "")
	harness.IsEqual(t, len(result.GetOpts("ddd")), 1, "")
	harness.IsNotNil(t, result.GetOpt("ddd"), "")
	harness.IsEqual(t, result.GetOpt("ddd").Optarg, "ARG", "")

	harness.IsTrue(t, result.HasOpt("d"), "")
	harness.IsEqual(t, len(result.GetOpts("d")), 1, "")
	harness.IsNotNil(t, result.GetOpt("d"), "")
	harness.IsEqual(t, result.GetOpt("d").Optarg, "ARG", "")

	args = split("-e123")
	result, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
	harness.IsTrue(t, result.HasOpt("eee"), "")
	harness.IsEqual(t, len(result.GetOpts("eee")), 1, "")
	harness.IsNotNil(t, result.GetOpt("eee"), "")
	harness.IsEqual(t, result.GetOpt("eee").Optarg, "123", "")

	harness.IsTrue(t, result.HasOpt("e"), "")
	harness.IsEqual(t, len(result.GetOpts("e")), 1, "")
	harness.IsNotNil(t, result.GetOpt("e"), "")
	harness.IsEqual(t, result.GetOpt("e").Optarg, "123", "")

	harness.IsFalse(t, result.HasOpt("@"), "")
	harness.IsEqual(t, len(result.GetOpts("@")), 0, "")
	harness.IsNil(t, result.GetOpt("@"), "")
}

func Test_ResultAPI_Negative(t *testing.T) {
	result := argp.ParseResult{
		Options: []argp.Result{},
		Args:    []string{},
	}
	harness.IsFalse(t, result.HasOpt("a"), "")
	harness.IsEqual(t, result.GetOpt("eee").WithDefault("DEFAULT"), "DEFAULT", "xxx")
	harness.IsNil(t, result.GetOpt("@"), "")
	harness.IsEqual(t, result.GetOpt("@").WithDefault("DEFAULT"), "DEFAULT", "")

	var res *argp.Result = nil
	harness.IsEqual(t, res.WithDefault("123"), "123", "nil call")
}

type testPairT1 struct {
	option argp.Option
	expect string
}

var helpCheckPatterns = []testPairT1{
	// ==========
	{
		option: argp.Option{Short: 'a', Long: "", Doc: "boolean option (short)"},
		expect: " -a                        boolean option (short)\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", Doc: "boolean option (long)"},
		expect: "     --aaa                 boolean option (long)\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", Doc: "boolean option"},
		expect: " -a, --aaa                 boolean option\n",
	},
	// ==========
	{
		option: argp.Option{Short: 'a', Long: "", ArgName: "ARG", Doc: "option with argument (short)"},
		expect: " -a ARG                    option with argument (short)\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", ArgName: "ARG", Doc: "option with argument (long)"},
		expect: "     --aaa ARG             option with argument (long)\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", ArgName: "ARG", Doc: "option with argument"},
		expect: " -a, --aaa ARG             option with argument\n",
	},
	// ==========
	{
		option: argp.Option{Short: 'a', Long: "", ArgName: "ARG", Flags: argp.OPTION_ARG_OPTIONAL, Doc: "option with argument (short)"},
		expect: " -a[ARG]                   option with argument (short)\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", ArgName: "ARG", Flags: argp.OPTION_ARG_OPTIONAL, Doc: "option with argument (long)"},
		expect: "     --aaa[=ARG]           option with argument (long)\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", ArgName: "ARG", Flags: argp.OPTION_ARG_OPTIONAL, Doc: "option with argument"},
		expect: " -a, --aaa[=ARG]           option with argument\n",
	},
	// ==========
	{
		option: argp.Option{Short: ' ', Long: " ", ArgName: " ", Flags: 0, Doc: "Documentation Line"},
		expect: "Documentation Line\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", ArgName: "ARG", Flags: argp.OPTION_HIDDEN, Doc: "hidden option"},
		expect: "",
	},
}

func testOption(t *testing.T, ptn *testPairT1) {

	buf := bytes.NewBufferString("")
	argp.PrintOptList(buf, []argp.Option{ptn.option})
	exp := strings.Split(ptn.expect, "\n")
	act := strings.Split(buf.String(), "\n")

	harness.IsEqual(t, len(act), len(exp), "")

	if len(exp) == len(act) {
		for i, _ := range exp {
			actLin := strings.TrimRight(act[i], " ")
			expLin := strings.TrimRight(exp[i], " ")

			harness.IsEqual(t, actLin, expLin, "")
		}
	}
}

func Test_OptionListPrinting(t *testing.T) {
	for _, pattern := range helpCheckPatterns {
		testOption(t, &pattern)
	}
}

func Test_OptionListPrinting2(t *testing.T) {
	option := []argp.Option{
		{Short: ' ', Long: "", ArgName: "", Doc: "OPTIONS:"},
		{Short: 'a', Long: "", ArgName: "", Doc: "enable option a"},
		{Short: 'b', Long: "b", ArgName: "", Doc: "run in mode b"},
		{Short: 's', Long: "silent", ArgName: "", Doc: "run in silent mode"},
		{Short: 'q', Long: "", ArgName: "", Flags: argp.OPTION_ALIAS, Doc: "this doc is ignored"},
		{Short: 'o', Long: "output", ArgName: "<file>", Doc: "specify the file to output"},
		{Short: '1', Long: "", ArgName: "", Doc: "run only once"},
		{Short: 0, Long: "", ArgName: "", Doc: ""},
		{Short: 0, Long: "", ArgName: "", Doc: "This line is for document"},
	}

	expect := "" +
		"OPTIONS:\n" +
		" -a                        enable option a\n" +
		" -b, --b                   run in mode b\n" +
		" -s, -q, --silent          run in silent mode\n" +
		" -o, --output <file>       specify the file to output\n" +
		" -1                        run only once\n" +
		"\n" +
		"This line is for document\n"

	buf := bytes.NewBufferString("")
	argp.PrintOptList(buf, option)
	exp := strings.Split(expect, "\n")
	act := strings.Split(buf.String(), "\n")

	harness.IsEqual(t, len(act), len(exp), "")

	if len(exp) == len(act) {
		for i, _ := range exp {
			actLin := strings.TrimRight(act[i], " ")
			expLin := strings.TrimRight(exp[i], " ")

			harness.IsEqual(t, actLin, expLin, "")
		}
	}
}

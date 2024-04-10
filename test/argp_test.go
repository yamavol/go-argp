package argp_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yamavol/go-argp"
	"github.com/yamavol/go-argp/test/harness"
)

var options = []argp.Option{
	{Short: ' ', Long: "", ArgName: "", Flags: 0, Doc: "CATEGORY 000:"},
	{Short: 'a', Long: "", ArgName: "", Flags: 0, Doc: "enable option a"},
	{Short: 'b', Long: "", ArgName: "", Flags: 0, Doc: "enable option b"},
	{Short: 'c', Long: "", ArgName: "", Flags: 0, Doc: "enable option c"},
	{Short: 'p', Long: "", ArgName: "", Flags: 0, Doc: "enable option p"},
	{Short: 'q', Long: "", ArgName: "", Flags: 0, Doc: "enable option q"},
	{Short: 'r', Long: "", ArgName: "", Flags: 0, Doc: "enable option r"},
	{Short: '1', Long: "", ArgName: "", Flags: 0, Doc: "enable option 1"},
	{Short: ' ', Long: "secret", ArgName: "", Flags: argp.OPTION_HIDDEN, Doc: "hidden option"},

	{Short: ' ', Long: "      ", ArgName: "     ", Flags: 0, Doc: "CATEGORY 111:"},
	{Short: 'o', Long: "output", ArgName: "<buf>", Flags: 0, Doc: "specify output buffer"},
	{Short: 'x', Long: "xxxx", ArgName: "<arg>", Flags: 0, Doc: "enable option x"},
	{Short: 'f', Long: "file", ArgName: "<file>", Flags: 0, Doc: "file to open"},
	{Short: ' ', Long: "ffff", ArgName: "", Flags: argp.OPTION_ALIAS, Doc: ""},
	{Short: ' ', Long: "fgfg", ArgName: "", Flags: argp.OPTION_ALIAS, Doc: ""},
	{Short: 'F', Long: "    ", ArgName: "", Flags: argp.OPTION_ALIAS, Doc: ""},
	{Short: 'K', Long: "kind", ArgName: "<kind>", Flags: argp.OPTION_ARG_OPTIONAL, Doc: "specify kind"},
	{Short: ' ', Long: "    ", ArgName: "", Flags: argp.OPTION_DOC, Doc: "kind is AAA BBB or CCC. Default is AAA"},

	{Short: ' ', Long: "    ", ArgName: "", Flags: 0, Doc: "CATEGORY 222:"},
	{Short: 'h', Long: "help", ArgName: "", Flags: 0, Doc: "print help and exit"},
	{Short: 'V', Long: "version", ArgName: "", Flags: 0, Doc: "print version and exit"},
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

	args = split("-x")
	_, err = argp.ParseArgs(options, args)
	harness.IsNotNil(t, err, "")

	args = split("--xxxx")
	_, err = argp.ParseArgs(options, args)
	harness.IsNotNil(t, err, "")

	args = split("--xxxx --")
	_, err = argp.ParseArgs(options, args)
	harness.IsNil(t, err, "")
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

type testPairT1 struct {
	option argp.Option
	expect string
}

var helpCheckPatterns = []testPairT1{
	// ==========
	{
		option: argp.Option{Short: 'a', Long: "", ArgName: "", Flags: 0, Doc: "boolean option (short)"},
		expect: " -a                        boolean option (short)\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", ArgName: "", Flags: 0, Doc: "boolean option (long)"},
		expect: "     --aaa                 boolean option (long)\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", ArgName: "", Flags: 0, Doc: "boolean option"},
		expect: " -a, --aaa                 boolean option\n",
	},
	// ==========
	{
		option: argp.Option{Short: 'a', Long: "", ArgName: "ARG", Flags: 0, Doc: "option with argument (short)"},
		expect: " -a ARG                    option with argument (short)\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", ArgName: "ARG", Flags: 0, Doc: "option with argument (long)"},
		expect: "     --aaa ARG             option with argument (long)\n",
	}, {
		option: argp.Option{Short: 'a', Long: "aaa", ArgName: "ARG", Flags: 0, Doc: "option with argument"},
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
		option: argp.Option{Short: ' ', Long: " ", ArgName: " ", Flags: 0, Doc: "CATEGORY HEADER:"},
		expect: "CATEGORY HEADER:\n",
	}, {
		option: argp.Option{Short: ' ', Long: "aaa", ArgName: "ARG", Flags: argp.OPTION_DOC, Doc: "Documentation Line"},
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
		{Short: ' ', Long: "", ArgName: "", Flags: 0, Doc: "OPTIONS:"},
		{Short: 'a', Long: "", ArgName: "", Flags: 0, Doc: "enable option a"},
		{Short: 'b', Long: "b", ArgName: "", Flags: 0, Doc: "run in mode b"},
		{Short: 's', Long: "silent", ArgName: "", Flags: 0, Doc: "run in silent mode"},
		{Short: 'q', Long: "", ArgName: "", Flags: argp.OPTION_ALIAS, Doc: "this doc is ignored"},
		{Short: 'o', Long: "output", ArgName: "<file>", Flags: 0, Doc: "specify the file to output"},
		{Short: '1', Long: "", ArgName: "", Flags: 0, Doc: "run only once"},
		{Short: 0, Long: "", ArgName: "", Flags: 0, Doc: ""},
		{Short: 0, Long: "", ArgName: "", Flags: argp.OPTION_DOC, Doc: "This line is for document"},
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

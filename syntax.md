
# Syntax

go-argp implements a option-syntax similar to the GNU argp library. The syntax is described in [this link](https://www.gnu.org/software/libc/manual/html_node/Argument-Syntax.html). This page shows the summary.

**short option rules**

    -a -b -c        ; Options begin with a hyphen delimitor.
    -abc            ; Multiple short options can be grouped together.
    -d ARG          ; Some option requires an argument. 
    -abcd ARG       ; Option with an argument can be grouped.
    -dARG           ; The space between the option and argument is optional. 
                    ;   If the ARG is "optional", this syntax must be used.
    --              ; Two-hyphen terminates all options. 
    -               ; Single hyphen is a non-optional argument.
    -d ABC -d DEF   ; Options may be supplied multiple times

**long option rules**

    --opt           ; Long option begins with two hyphen delimitor
    --opt ARG       ; Some option requires an argument. 
    --opt=ARG       ; = sign can be used as a separator.
                    ;   If the ARG is optional, this syntax must be used
                    ;   because it is ambiguous.
    --oo AB --oo YZ ; Options may be supplied multiple times.

**other option rules**

    ARG0 ARG1 -xyz  ; Non-option can appear before the options. This is against
                    ;   the POSIX standard.

**unsupported syntax**

    -I file1 file2  ; Multi-argument option is not supported. It is ambiguous.
    -a=ARG          ; = cannot be used with short option. It is ambigous.

**go flag's rules**

    -opt --opt      ; Flag does not distinguish short and long names
    -opt ARG        ; The Flag accepts an argument
    -opt=ARG        ; = sign can be used as a separator
                    ;
    -opt            ; one option can be passed only once



    















build-in help (-h and --help) message
scan arguments into struct fields with configuration in tags
scan into composite field types (arrays, slices, structs)
allow for nested sub commands
GNU command line argument rules:

arguments are options when they begin with a hyphen -
multiple options can be combined: -abc is the same as -a -b -c
long options start with two hyphens: --abc is one option
option names are alphanumeric characters
options can have a value: -a 1 means that a has value 1
option values can be separated by a space, equal sign, or nothing: -a1 -a=1 -a 1 are all equal
options and non-options can be interleaved
the argument -- terminates all options so that all following arguments are treated as non-options
a single - argument is a non-option usually used to mean standard in or out streams
options may be specified multiple times, only the last one determines its value
options can have multiple values: -a 1 2 3 means that a is an array/slice/struct of three numbers of value [1,2,3]







## FORMAT WIDTH

    | -a, --longname <param> |
    |--->                    | reserve 4 letters for the short option
    |----------------------->| reserve 25 letters for all options

## FORMAT PATTERNS

no arg:

    -o               
    -o, --opt
        --opt

with arg:

    -o ARG           
    -o, --opt ARG
        --opt ARG

with optional-arg:

    -o[ARG]          
    -o, --opt[=ARG]
        --opt[=ARG]

with alias:

    -o, -p		
    -o, -p, --opt, --ppt
        --opt, --ppt

with alias and arg:

    -o ARG, -p ARG		
    -o, -p, --opt ARG,
        --opt ARG, --ppt ARG

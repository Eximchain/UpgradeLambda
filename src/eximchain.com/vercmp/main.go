package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/go-version"
)

var (
	debug      bool
	noexitcode bool
)

func compareV(v1S, v2S string) (result int) {
	v1, _ := version.NewVersion(v1S)
	v2, _ := version.NewVersion(v2S)
	if v1.LessThan(v2) {
		result = -1
	} else if v1.Equal(v2) {
		result = 0
	} else if v1.GreaterThan(v2) {
		result = 1
	}
	return
}

func compare(regularexpression, inputStr, compareStr string, lessthan, equals, morethan, errcode int) (result int) {
	if regularexpression == "" || compareStr == "" {
		if debug {
			fmt.Printf("regex: '%+v' or compare: '%s' is empty!\n", regularexpression, compareStr)
		}
		result = errcode
		return
	}
	regex, err := regexp.Compile(regularexpression)
	if err != nil {
		if debug {
			fmt.Printf("Error: %+v", err)
		}
		result = errcode
	} else {
		regexResult := regex.Find([]byte(inputStr))
		v1S := string(regexResult)

		compareResult := regex.Find([]byte(compareStr))
		v2S := string(compareResult)
		if v2S == "" {
			v2S = compareStr
		}

		v1, _ := version.NewVersion(v1S)
		v2, _ := version.NewVersion(v2S)
		if v1 == nil || v2 == nil {
			if debug {
				fmt.Printf("v1: %v, or v2: %v is nil\n", v1, v2)
			}
			result = errcode
		} else {
			var s string
			switch v1.Compare(v2) {
			case -1:
				s = "v1: %v < v2: %v\n"
				result = lessthan
			case 0:
				s = "v1: %v = v2: %v\n"
				result = equals
			case 1:
				s = "v1: %v > v2: %v\n"
				result = morethan
			}
			if debug {
				fmt.Printf(s, v1, v2)
			}
			if noexitcode {
				fmt.Print(result)
				result = 0
			}
		}
	}
	return
}

func max(a, b int) (result int) {
	if a > b {
		result = a
	} else {
		result = b
	}
	return
}

func main() {
	var (
		lessthan, morethan, equal, errcode int
		maxval                             int
		regularexpression, compareStr      string
		inputStr                           string
	)
	flag.IntVar(&lessthan, "lessthan", -1, "Specifies the return value for less than")
	flag.IntVar(&equal, "equals", 0, "Specifies the return value for equal result")
	flag.IntVar(&morethan, "morethan", 1, "Specifies the return value for more than")
	flag.IntVar(&errcode, "error", 4, "Error parsing the regular expression")
	flag.BoolVar(&debug, "debug", false, "Shows verbose output")
	flag.BoolVar(&noexitcode, "noexitcode", false, "Returns result to stdout, and exit code as 0 (when true)")
	flag.StringVar(&regularexpression, "regex", `(([0-9]+)\.?)+`, "Regular expression to match")
	flag.StringVar(&compareStr, "compare", "", "the string to compare the input against.")
	flag.StringVar(&inputStr, "input", "", "the input (optional. If not present, reads from stdin)")
	if len(os.Args) < 1 {
		fmt.Println("Eximchain version comparison utility")
		fmt.Println()
		fmt.Println(`Finds the regex string from stdin and compare it against the "compare" string`)
		fmt.Println("If the parsed input < compare, return the value specified in lessthan.")
		fmt.Println("If the parsed input = compare, return the value specified in equals.")
		fmt.Println("If the parsed input > compare, return the value specified in morethan.")
		fmt.Println("If the regular expression (regex) cannot be parsed, or no input or regularexpression is provided, return the value specified in error.")
		flag.Usage()
		return
	}
	flag.Parse()

	// ensures that error code is always higher than lessthan, equals, or morethan
	maxval = max(lessthan, equal)
	maxval = max(morethan, maxval)
	if maxval >= errcode {
		errcode = maxval + 1
	}

	if inputStr == "" {
		reader := bufio.NewReader(os.Stdin)
		inputStr, _ = reader.ReadString('\n')
	}
	if debug {
		fmt.Printf("Read/Input string: '%s'\n", inputStr)
	}

	os.Exit(compare(regularexpression, inputStr, compareStr, lessthan, equal, morethan, errcode))

}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ffuf/ffuf/pkg/ffuf"
)

type UsageSection struct {
	Name          string
	Description   string
	Flags         []UsageFlag
	Hidden        bool
	ExpectedFlags []string
}

//PrintSection prints out the section name, description and each of the flags
func (u *UsageSection) PrintSection(max_length int, extended bool) {
	// Do not print if extended usage not requested and section marked as hidden
	if !extended && u.Hidden {
		return
	}
	fmt.Printf("%s:\n", u.Name)
	for _, f := range u.Flags {
		f.PrintFlag(max_length)
	}
	fmt.Printf("\n")
}

type UsageFlag struct {
	Name        string
	Description string
	Default     string
}

//PrintFlag prints out the flag name, usage string and default value
func (f *UsageFlag) PrintFlag(max_length int) {
	// Create format string, used for padding
	format := fmt.Sprintf("  -%%-%ds %%s", max_length)
	if f.Default != "" {
		format = format + " (default: %s)\n"
		fmt.Printf(format, f.Name, f.Description, f.Default)
	} else {
		format = format + "\n"
		fmt.Printf(format, f.Name, f.Description)
	}
}

func Usage() {
	u_http := UsageSection{
		Name:          "HTTP OPTIONS",
		Description:   "Options controlling the HTTP request and its parts.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"H", "X", "b", "d", "r", "u", "recursion", "recursion-depth", "replay-proxy", "timeout", "ignore-body", "x"},
	}
	u_general := UsageSection{
		Name:          "GENERAL OPTIONS",
		Description:   "",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"ac", "acc", "c", "maxtime", "maxtime-job", "p", "rate", "s", "sa", "se", "sf", "t", "v", "V"},
	}
	u_compat := UsageSection{
		Name:          "COMPATIBILITY OPTIONS",
		Description:   "Options to ensure compatibility with other pieces of software.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        true,
		ExpectedFlags: []string{"compressed", "cookie", "data", "data-ascii", "data-binary", "i", "k"},
	}
	u_matcher := UsageSection{
		Name:          "MATCHER OPTIONS",
		Description:   "Matchers for the response filtering.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"mc", "ml", "mr", "ms", "mw"},
	}
	u_filter := UsageSection{
		Name:          "FILTER OPTIONS",
		Description:   "Filters for the response filtering.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"fc", "fl", "fr", "fs", "fw"},
	}
	u_input := UsageSection{
		Name:          "INPUT OPTIONS",
		Description:   "Options for input data for fuzzing. Wordlists and input generators.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"D", "ic", "input-cmd", "input-num", "mode", "request", "request-proto", "e", "w"},
	}
	u_output := UsageSection{
		Name:          "OUTPUT OPTIONS",
		Description:   "Options for output. Output file formats, file names and debug file locations.",
		Flags:         make([]UsageFlag, 0),
		Hidden:        false,
		ExpectedFlags: []string{"debug-log", "o", "of", "od"},
	}
	sections := []UsageSection{u_http, u_general, u_compat, u_matcher, u_filter, u_input, u_output}

	// Populate the flag sections
	max_length := 0
	flag.VisitAll(func(f *flag.Flag) {
		found := false
		for i, section := range sections {
			if strInSlice(f.Name, section.ExpectedFlags) {
				sections[i].Flags = append(sections[i].Flags, UsageFlag{
					Name:        f.Name,
					Description: f.Usage,
					Default:     f.DefValue,
				})
				found = true
			}
		}
		if !found {
			fmt.Printf("DEBUG: Flag %s was found but not defined in help.go.\n", f.Name)
			os.Exit(1)
		}
		if len(f.Name) > max_length {
			max_length = len(f.Name)
		}
	})

	fmt.Printf("Fuzz Faster U Fool - v%s\n\n", ffuf.VERSION)

	// Print out the sections
	for _, section := range sections {
		section.PrintSection(max_length, false)
	}

}

func strInSlice(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

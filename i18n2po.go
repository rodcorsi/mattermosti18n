package main

import (
	"./lib"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var template bool
var output string

func init() {
	flag.BoolVar(&template, "pot", false, "Define if the output is a template (POT)")
	flag.BoolVar(&template, "p", false, "Define if the output is a template (POT)")

	flag.StringVar(&output, "output", "", "Define the output file name")
	flag.StringVar(&output, "o", "", "Define the output file name")
}

func loadFile(filename string) []byte {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not load ", filename, err)
		os.Exit(1)
	}
	return file
}

func main() {
	var filename string
	var targetfilename string

	// Get input, if any...
	flag.Parse()
	if flag.NArg() == 1 {
		filename = flag.Arg(0)

	} else if flag.NArg() == 2 {
		filename = flag.Arg(0)
		targetfilename = flag.Arg(1)

	} else {
		fmt.Println(usage)
		os.Exit(1)
	}

	var s, t *mattermosti18n.Translations

	s = mattermosti18n.LoadJson(loadFile(filename))
	if len(targetfilename) > 0 {
		t = mattermosti18n.LoadJson(loadFile(targetfilename))
	}

	if len(output) == 0 {
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		if template {
			output = name + "_new.pot"
		} else {
			output = name + "_new.po"
		}
	}

	err := ioutil.WriteFile(output, s.ToPO(t, template), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const usage = `i18n2po [Options] [source filename] [target filename]
Convert json translation to PO/POT files

[source filename]
    Json file in default language normaly en

[target filename]
    Translated json file

Options:
    -p, -pot       Define if the output is a template (POT)
    -o, -output    Define the output file name`

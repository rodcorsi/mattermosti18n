package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/rodcorsi/mattermosti18n"
)

var pot bool
var output string
var template string

func init() {
	flag.BoolVar(&pot, "pot", false, "Define if the output is a template (POT)")
	flag.BoolVar(&pot, "p", false, "Define if the output is a template (POT)")

	flag.StringVar(&output, "output", "", "Define the output file name")
	flag.StringVar(&output, "o", "", "Define the output file name")

	flag.StringVar(&template, "template", "", "PO file that will be filled with Json translations")
	flag.StringVar(&template, "t", "", "PO file that will be filled with Json translations")
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

	// Get input, if any...
	flag.Parse()
	if flag.NArg() == 1 {
		filename = flag.Arg(0)
	} else {
		fmt.Println(usage)
		os.Exit(1)
	}

	var json, po *mattermosti18n.Translations

	json = mattermosti18n.LoadJson(loadFile(filename))
	if len(template) > 0 {
		po = mattermosti18n.LoadPO(loadFile(template))
	}

	if len(output) == 0 {
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		if pot {
			output = name + "_new.pot"
		} else {
			output = name + "_new.po"
		}
	}

	var err error
	if len(template) > 0 {
		err = ioutil.WriteFile(output, po.ToPO(json, pot), 0644)
	} else {
		err = ioutil.WriteFile(output, json.ToPO(nil, pot), 0644)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const usage = `i18n2po [Options] [target filename]
Convert json translation to PO/POT files

[target filename]
    Translated Json file

Options:
    -t, -template  PO file that will be filled with Json translations
    -p, -pot       Define if the output is a template (POT)
    -o, -output    Define the output file name (PO/POT)`

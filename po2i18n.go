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

var template string
var output string

func init() {
	flag.StringVar(&template, "template", "", "Define the template to output in the same order")
	flag.StringVar(&template, "t", "", "Define the template to output in the same order")

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

	// Get input, if any...
	flag.Parse()
	if flag.NArg() == 1 {
		filename = flag.Arg(0)
	} else {
		fmt.Println(usage)
		os.Exit(1)
	}

	po := mattermosti18n.LoadPO(loadFile(filename))

	var json *mattermosti18n.Translations
	if len(template) > 0 {
		json = mattermosti18n.LoadJson(loadFile(template))
	}

	if len(output) == 0 {
		name := strings.TrimSuffix(filename, filepath.Ext(filename))
		output = name + "_new.json"
	}

	err := ioutil.WriteFile(output, po.ToJson(json), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const usage = `po2i18n [Options] filename
Convert PO translation to json files

[filename]
    PO file

Options:
    -t, -template   Define the json template file to output in the same order
    -o, -output     Define the output file name`

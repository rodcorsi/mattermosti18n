package mattermosti18n

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
[
	{ "id": "xxx", "translation": "xxxx" },
	{ "id": "xxx", "translation": "xxxx" }
]
*/
type Platform []struct {
	Id          string `json:"id"`
	Translation string `json:"translation"`
}

type PO struct {
	Localization string //#:
	Original     string //msgid
	Translation  string //msgstr
}

type Translations struct {
	Webstatic bool
	Header    string
	Order     []string
	Data      map[string]PO //[msgctxt]
}

func unquote(text string) string {
	s, err := strconv.Unquote(text)
	if err != nil {
		return text
	}
	return s
}

func parseWebStaticJson(data []byte, trans *Translations) error {
	var parse map[string]string

	err := json.Unmarshal(data, &parse)
	if err != nil {
		return err
	}

	(*trans).Webstatic = true
	(*trans).Header = header

	(*trans).Data = make(map[string]PO)
	for k, t := range parse {
		(*trans).Data[k] = PO{Localization: "." + k, Original: t, Translation: t}
	}

	//Create original order
	//Get "xx":""
	(*trans).Order = regexp.MustCompile(`"([^"\\\\]*|\\\\["\\\\bfnrt\/]|\\\\u[0-9a-f]{4}|\\")*"\s*:`).FindAllString(string(data), -1)
	reg := regexp.MustCompile(`"\s*:\s*$`)

	for i, w := range (*trans).Order {
		ss := reg.ReplaceAllString(w, "\"")
		(*trans).Order[i] = unquote(ss)
	}

	return nil
}

func parsePlatformJson(data []byte, trans *Translations) error {
	var parse Platform
	err := json.Unmarshal(data, &parse)

	if err != nil {
		return err
	}

	(*trans).Webstatic = false
	(*trans).Header = header

	(*trans).Data = make(map[string]PO)
	for _, t := range parse {
		(*trans).Data[t.Id] = PO{Localization: "." + t.Id, Original: t.Translation, Translation: t.Translation}
	}

	//Create original order
	//Get "id":"xxxxxxx"
	(*trans).Order = regexp.MustCompile(`\s*"id"\s*:\s*"([^"\\\\]*|\\\\["\\\\bfnrt\/]|\\\\u[0-9a-f]{4}|\\")*"`).FindAllString(string(data), -1)
	reg := regexp.MustCompile(`\s*"id"\s*:\s*`)

	for i, w := range (*trans).Order {
		ss := reg.ReplaceAllString(w, "")
		(*trans).Order[i] = unquote(ss)
	}

	return nil
}

func LoadJson(data []byte) *Translations {
	var parse Translations

	err := parseWebStaticJson(data, &parse)
	if err != nil {
		//try to parse Platform type
		err2 := parsePlatformJson(data, &parse)
		if err2 != nil {
			fmt.Println("Could not parse ", err, err2)
			os.Exit(1)
		}
	}
	return &parse
}

func (source *Translations) ToPO(target *Translations, template bool) []byte {
	var buf bytes.Buffer
	buf.WriteString((*source).Header)

	hastarget := target != nil && len((*target).Data) > 0

	var k, t string
	var po PO
	for i := 0; i < len((*source).Order); i = i + 1 {
		k = (*source).Order[i]
		po = (*source).Data[k]

		if hastarget {
			t = strconv.Quote((*target).Data[k].Translation) //translation in target language
		} else {
			t = strconv.Quote(po.Translation) //translation in source language (en)
		}

		buf.WriteString(fmt.Sprintln())
		buf.WriteString(fmt.Sprintf("#: %v\n", po.Localization))
		buf.WriteString(fmt.Sprintln("msgctxt", strconv.Quote(k)))
		buf.WriteString(fmt.Sprintln("msgid", strconv.Quote(po.Original)))

		if template {
			buf.WriteString(fmt.Sprintln("msgstr", `""`))
		} else {
			buf.WriteString(fmt.Sprintln("msgstr", t))
		}
	}

	return buf.Bytes()
}

func readField(text string, scanner *bufio.Scanner) (string, bool) {
	ret := unquote(text)

	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(ln, "\"") {
			ret = ret + unquote(ln)
		} else {
			return ret, true
		}
	}
	return ret, false
}

func LoadPO(data []byte) *Translations {
	var fields []string
	var id, t, local, original string

	reg := regexp.MustCompile(" +")

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var parse Translations
	parse.Data = make(map[string]PO)
	parse.Order = make([]string, 0)

	//Read header
	parse.Header = ""
	next := scanner.Scan()
	for next {
		ln := strings.TrimSpace(scanner.Text())
		if len(ln) == 0 {
			break
		}
		parse.Header = parse.Header + ln + "\n"
		next = scanner.Scan()
	}

	//Read translations
	for next {
		ln := strings.TrimSpace(scanner.Text())

		if len(ln) == 0 || strings.HasPrefix(ln, "\"") {
			next = scanner.Scan()
			continue
		}

		fields = reg.Split(ln, 2)
		if len(fields) != 2 {
			continue
		}

		switch fields[0] {
		case "#:":
			local, next = readField(fields[1], scanner)
			continue
		case "msgctxt":
			id, next = readField(fields[1], scanner)
			parse.Order = append(parse.Order, id)
			continue
		case "msgid":
			original, next = readField(fields[1], scanner)
			continue
		case "msgstr":
			if len(id) > 0 {
				t, next = readField(fields[1], scanner)
				parse.Data[id] = PO{Localization: local, Original: original, Translation: t}
				continue
			}
		}
		next = scanner.Scan()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &parse
}

func (source *Translations) toJsonWebStatic(template *Translations) []byte {
	var buf bytes.Buffer
	var k, t string
	var order []string
	next := false

	hastemplate := template != nil && len((*template).Order) > 0

	if hastemplate {
		order = (*template).Order
	} else {
		order = (*source).Order
	}

	buf.WriteString("{\n")

	for i := 0; i < len(order); i++ {
		k = order[i]
		t = (*source).Data[k].Translation
		if len(t) == 0 && hastemplate {
			t = (*template).Data[k].Translation
		}

		if next {
			buf.WriteString(",\n")
		} else {
			next = true
		}
		buf.WriteString(fmt.Sprintf("%v%v: %v", indent, strconv.Quote(k), strconv.Quote(t)))
	}
	buf.WriteString("\n}")
	return buf.Bytes()
}

func (source *Translations) toJsonPlatform(template *Translations) []byte {
	var buf bytes.Buffer
	var k, t string
	var order []string
	next := false

	hastemplate := template != nil && len((*template).Order) > 0

	if hastemplate {
		order = (*template).Order
	} else {
		order = (*source).Order
	}

	buf.WriteString("[\n")

	for i := 0; i < len(order); i++ {
		k = order[i]
		t = (*source).Data[k].Translation
		if len(t) == 0 && hastemplate {
			t = (*template).Data[k].Translation
		}

		if next {
			buf.WriteString(",\n")
		} else {
			next = true
		}
		buf.WriteString(fmt.Sprintf("%v{\n", indent))
		buf.WriteString(fmt.Sprintf("%v\"id\": %v,\n", indent2x, strconv.Quote(k)))
		buf.WriteString(fmt.Sprintf("%v\"translation\": %v\n", indent2x, strconv.Quote(t)))
		buf.WriteString(fmt.Sprintf("%v}", indent))
	}
	buf.WriteString("\n]")
	return buf.Bytes()
}

func (source *Translations) ToJson(template *Translations) []byte {
	if template != nil && (*template).Webstatic {
		return source.toJsonWebStatic(template)
	}
	return source.toJsonPlatform(template)
}

const indent = "  "
const indent2x = indent + indent

const header = `# MATTERMOST.` + "\n" +
	`# Copyright (C) Mattermost` + "\n" +
	`# This file is distributed under the same license as the PACKAGE package.` + "\n" +
	`# FIRST AUTHOR <EMAIL@ADDRESS>, YEAR.` + "\n" +
	`#` + "\n" +
	`#, fuzzy` + "\n" +
	`msgid ""` + "\n" +
	`msgstr ""` + "\n" +
	`"Project-Id-Version: Translate mattermost 1.0\n"` + "\n" +
	`"Report-Msgid-Bugs-To: translate@mattermost.org\n"` + "\n" +
	`"POT-Creation-Date: 2016-02-16 10:03+0200\n"` + "\n" +
	`"PO-Revision-Date: YEAR-MO-DA HO:MI+ZONE\n"` + "\n" +
	`"Last-Translator:\n"` + "\n" +
	`"Language-Team:\n"` + "\n" +
	`"MIME-Version: 1.0\n"` + "\n" +
	`"Content-Type: text/plain; charset=UTF-8\n"` + "\n" +
	`"Content-Transfer-Encoding: 8bit\n"` + "\n" +
	`"X-Generator: i18n2po\n"` + "\n"

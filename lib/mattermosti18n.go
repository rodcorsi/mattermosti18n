package mattermosti18n

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"strconv"
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

/*
{
	"id":"translation",
	"id":"translation"
}
*/
type Translations struct {
	Webstatic bool
	Order     []string
	Data      map[string]string
}

func unquote(text string) string {
	s, err := strconv.Unquote(text)
	if err != nil {
		fmt.Println("Error on Unquote ", err)
		return text
	}
	return s
}

func appendquoted(text1 string, text2 string) string {
	return text1 + unquote(text2)
}

func parseWebStaticJson(data []byte, trans *Translations) error {
	err := json.Unmarshal(data, &(*trans).Data)
	if err != nil {
		return err
	}

	(*trans).Webstatic = true
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

	(*trans).Data = make(map[string]string)
	for _, t := range parse {
		(*trans).Data[t.Id] = t.Translation
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
	buf.WriteString(header)

	notarget := target == nil || len((*target).Data) == 0

	var fixed, k, t string
	for i := 0; i < len((*source).Order); i = i + 1 {
		k = (*source).Order[i]
		t = (*source).Data[k]

		if notarget {
			fixed = strconv.Quote(t)
		} else {
			fixed = strconv.Quote((*target).Data[k])  //translation in source language (en)
		}

		buf.WriteString(fmt.Sprintln())
		buf.WriteString(fmt.Sprintf("#: .%v\n", k))
		buf.WriteString(fmt.Sprintln("msgctxt", strconv.Quote(k)))
		buf.WriteString(fmt.Sprintln("msgid", fixed))

		if template {
			buf.WriteString(fmt.Sprintln("msgstr", `""`))
		} else {
			buf.WriteString(fmt.Sprintln("msgstr", fixed))
		}
	}

	return buf.Bytes()
}

func LoadPO(data []byte) *Translations {
	var fields []string
	var id, translation string

	reg := regexp.MustCompile(" +")

	scanner := bufio.NewScanner(bytes.NewReader(data))
	var parse Translations
	parse.Data = make(map[string]string)
	parse.Order = make([]string, 0)

	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())

		if len(ln) == 0 || strings.HasPrefix(ln, "#") {
			continue
		}

		if strings.HasPrefix(ln, "\"") {
			if len(id) == 0 {
				continue
			}
			translation = appendquoted(translation, ln)
			parse.Data[id] = translation
			continue
		}

		fields = reg.Split(ln, 2)
		if len(fields) != 2 {
			fmt.Println("Error on Split fields", fields)
			os.Exit(1)
		}

		if fields[0] == "msgctxt" {
			id = unquote(fields[1])
			parse.Order = append(parse.Order, id)
			translation = ""
		} else if fields[0] == "msgstr" {
			if len(id) == 0 {
				continue
			}
			translation = appendquoted(translation, fields[1])
			parse.Data[id] = translation
		}
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
	next := false

	buf.WriteString("{\n")

	for i := 0; i < len((*template).Order); i++ {
		k = (*template).Order[i]
		t = (*source).Data[k]
		if len(t) == 0 {
			t = (*template).Data[k]
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
	next := false

	buf.WriteString("[\n")

	for i := 0; i < len((*template).Order); i++ {
		k = (*template).Order[i]
		t = (*source).Data[k]
		if len(t) == 0 {
			t = (*template).Data[k]
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
	if (*template).Webstatic {
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

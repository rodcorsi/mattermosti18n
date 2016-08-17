package mattermosti18n

import (
	"testing"
)

const webstatic = `{
  "about.teamEdtion": "Team \n\"Edition\""
}`

const platform = `[
  {
  	"id": "about.teamEdtion",
  	"translation": "Team \n\"Edition\""
  }
]`

const po = header +
	`
#: .about.teamEdtion
msgctxt "about.teamEdtion"
msgid "Team \n\"Edition\""
msgstr "Team \n\"Edition\""
`

const potrans = header +
	`
#: .about.teamEdtion
msgctxt "about.teamEdtion"
msgid "Team \n\"Edition\""
msgstr "Equipe \n\"Edição\""
`

func TestLoad(t *testing.T) {
	expect := "Team \n\"Edition\""

	m := LoadJson([]byte(webstatic))
	result := (*m).Data["about.teamEdtion"].Translation
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m = LoadJson([]byte(platform))
	result = (*m).Data["about.teamEdtion"].Translation
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m = LoadPO([]byte(po))
	result = (*m).Data["about.teamEdtion"].Translation
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}
}

const poPlural = header +
	`
#: .about.teamEdtion
msgctxt "about.teamEdtion"
msgid "One"
msgid_plural "Other"
msgstr[0] "One"
msgstr[1] "Other"
`

func TestLoadPlural(t *testing.T) {
	m := LoadPO([]byte(poPlural))
	result := (*m).Data["about.teamEdtion"].Plural
	if result != "Other" {
		t.Fatal("Expect: Other\nResult:", result)
	}

	result = (*m).Data["about.teamEdtion"].TransPlural[0]
	if result != "One" {
		t.Fatal("Expect: One\nResult:", result)
	}

	result = (*m).Data["about.teamEdtion"].TransPlural[1]
	if result != "Other" {
		t.Fatal("Expect: Other\nResult:", result)
	}
}

func TestToPO(t *testing.T) {
	expect := po

	m := LoadJson([]byte(webstatic))
	result := string(m.ToPO(nil, false))
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m = LoadJson([]byte(platform))
	result = string(m.ToPO(nil, false))
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m2 := LoadPO([]byte(po))
	result = string(m2.ToPO(m, false))
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m = LoadPO([]byte(po))
	result = string(m.ToPO(nil, false))
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}

	m = LoadPO([]byte(poPlural))
	expect = poPlural
	result = string(m.ToPO(nil, false))
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}
}

const po_lines = `
""
msgctxt "about\n."
"team\"E\"dtion"
any "jsjsj"
msgstr ""
"isso e um grande texto"
" mesmos"
`

func TestPOMultilines(t *testing.T) {
	m := LoadPO([]byte(po_lines))

	key := "about\n.team\"E\"dtion"
	expect := "isso e um grande texto mesmos"
	result := m.Data[key].Translation
	if result != expect {
		t.Fatal("Expect:", expect, "\nResult:", result)
	}
}

func TestPOTranslated(t *testing.T) {
	m := LoadPO([]byte(po))
	target := LoadPO([]byte(potrans))
	result := string(m.ToPO(target, false))
	if result != potrans {
		t.Fatal("Expect:", potrans, "\nResult:", result)
	}
}

const json_plural = `[
  {
    "id": "id1",
    "translation": {
      "one": "aaa",
      "other": "bbb"
    }
  }
]
`

func TestParseJSONPlural(t *testing.T) {
	m := LoadJson([]byte(json_plural))

	if m == nil {
		t.Fatal("Unable to Parse Plural JSON")
	}

	result := string(m.ToJson(nil))
	if result != json_plural {
		t.Fatal("Expect:", json_plural, "\nResult:", result)
	}
}

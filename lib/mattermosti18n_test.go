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

func TestLoad(t *testing.T) {
	expect := "Team \n\"Edition\""

	m := LoadJson([]byte(webstatic))
	result := (*m).Data["about.teamEdtion"]
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	} 

	m = LoadJson([]byte(platform))
	result = (*m).Data["about.teamEdtion"]
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	} 

	m = LoadPO([]byte(po))
	result = (*m).Data["about.teamEdtion"]
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	} 
}

func TestToPO(t *testing.T) {
	expect := po

	m := LoadJson([]byte(webstatic))
	result := string(m.ToPO(nil,false))
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	}

	m = LoadJson([]byte(platform))
	result = string(m.ToPO(nil,false))
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	}

	m = LoadPO([]byte(po))
	result = string(m.ToPO(m,false))
	if result != expect {
		t.Error("Expect:", expect, "\nResult:", result)
	}
}

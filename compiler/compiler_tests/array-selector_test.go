package compiler_tests

import (
	"encoding/json"
	"testing"

	cp "me.compiler/compiler"
)

func TestArraySelector(t *testing.T) {
	jsonData := `{
		"a": {
			"b": "AB",
			"c": [
				{
					"c": {
						"c1": 1,
						"c2": "C"
					},
					"d": {
						"e": {
							"f": "F"
						}
					}
				},
				{
					"c": {
						"c1": 2,
						"c2": "X"
					},
					"d": {
						"e": {
							"f": "Y"
						}
					}
				}
			]
		}
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		SetTxt("data.a.c[c.c1===1].d.e.f")
	s, _, _ := c.Compile().GetResponse()
	if s != "F" {
		t.Errorf("1) compile %s expected %s", s, "F")
	}

	c.SetTxt("data.a.c[c.c1===2].d.e.f")
	s, _, _ = c.Compile().GetResponse()
	if s != "Y" {
		t.Errorf("2) compile %s expected %s", s, "Y")
	}

	c.SetTxt("data.a.c[c.c2===C].d.e.f")
	s, _, _ = c.Compile().GetResponse()
	if s != "F" {
		t.Errorf("3) compile %s expected %s", s, "F")
	}

	c.SetTxt("data.a.c[c.c2===X].d.e.f")
	s, _, _ = c.Compile().GetResponse()
	if s != "Y" {
		t.Errorf("4) compile %s expected %s", s, "Y")
	}

	c.SetTxt("data.a.b[d===4].e")
	_, err := c.Compile().IsOk()
	if err == nil {
		t.Errorf("5) compile error, error expected")
	}

	c.SetTxt("data.a.c[c.c2===W].d.e.f")
	_, err = c.Compile().IsOk()
	if err == nil {
		t.Errorf("6) compile error, error expected")
	}
}

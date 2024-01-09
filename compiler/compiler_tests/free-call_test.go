package compiler_tests

import (
	"encoding/json"
	"testing"

	cp "me.compiler/compiler"
)

func TestFreeCall(t *testing.T) {
	jsonData := `{
		"a": {
			"b": "AB",
			"c": "AC"
		}
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		SetTxt("data.a.b")

	s, _, _ := c.Compile().GetResponse()
	if s != "AB" {
		t.Errorf("1) compile %s expected %s", s, "AB")
	}

	c.SetTxt("data.a.c")
	s, _, _ = c.Compile().GetResponse()
	if s != "AC" {
		t.Errorf("2) compile %s expected %s", s, "AC")
	}

	c.SetTxt("data.a.d")
	_, err := c.Compile().IsOk()
	if err == nil {
		t.Errorf("3) compile error, error expected")
	}
}

func TestFreeCallArrayIndex(t *testing.T) {
	jsonData := `{
		"a": {
			"b": "AB",
			"c": [
				{
					"c1": "AB0C1"
				},
				{
					"c2": "AB1C2"
				}
			]
		}
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		SetTxt("data.a.c.0.c1")

	s, _, _ := c.Compile().GetResponse()
	if s != "AB0C1" {
		t.Errorf("1) compile %s expected %s", s, "AB0C1")
	}
}

package compiler_tests

import (
	"encoding/json"
	"testing"

	cp "me.compiler/compiler"
)

func TestCountLabel(t *testing.T) {
	jsonData := `{
		"c": [
			{
				"c": 1
			},
			{
				"c": 2
			}
		],
		"d": [
			{
				"d": 1
			}
		],
		"e": []
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		NewSetAuxArrayAddOn("data.c").
		RunPreAddOns().
		SetTxt("[count|one|two]")

	s, _, _ := c.Compile().GetResponse()
	if s != "2 two" {
		t.Errorf("1) compile %s expected %s", s, "2 two")
	}

	c.
		Reset(cp.PreAddOns).
		NewSetAuxArrayAddOn("data.d").
		RunPreAddOns().
		Reset(cp.Parser)
	s, _, _ = c.
		Compile().
		GetResponse()
	if s != "1 one" {
		t.Errorf("2) compile %s expected %s", s, "1 one")
	}

	c.
		Reset(cp.PreAddOns).
		NewSetAuxArrayAddOn("data.e").
		RunPreAddOns().
		Reset(cp.Parser)
	_, err := c.Compile().IsOk()
	if err == nil {
		t.Errorf("3) compile error, error expected")
	}
}

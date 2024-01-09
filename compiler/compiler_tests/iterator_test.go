package compiler_tests

import (
	"encoding/json"
	"testing"

	cp "me.compiler/compiler"
)

func TestIterator(t *testing.T) {
	jsonData := `{
		"c": [
			{
				"c": "1"
			},
			{
				"c": "2"
			}
		],
		"d": [
			{
				"d": "1"
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
		SetTxt("data.++.c")

	expected := make([]string, 0)
	expected = append(expected, "1", "2")

	_, r, b := c.Compile().GetResponse()
	for i, v := range r {
		if v != expected[i] {
			t.Errorf("1.%d) compile %v expected %v", i, r, expected[i])
		}
	}
	if !b {
		t.Errorf("1.-) should tell is multiparts")
	}

	c.
		Reset(cp.PreAddOns).
		NewSetAuxArrayAddOn("data.d").
		RunPreAddOns().
		SetTxt("data.++.d")

	expected = make([]string, 0)
	expected = append(expected, "1")

	_, r, b = c.Compile().GetResponse()
	for i, v := range r {
		if v != expected[i] {
			t.Errorf("2.%d) compile %v expected %v", i, r, expected[i])
		}
	}
	if !b {
		t.Errorf("2.-) should tell is multiparts")
	}
}

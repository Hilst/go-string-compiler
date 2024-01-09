package compiler_tests

import (
	"encoding/json"
	"strings"
	"testing"

	cp "me.compiler/compiler"
)

func TestComplexConcat(t *testing.T) {
	jsonData := `{
		"a": {
			"b": "AB",
			"c": [
				{
					"d": 1,
					"e": "AC0E"
				}
			]
		}
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		SetTxt("before 1 @{data.a.b} middle 2 @{data.a.c[d===1].e} last 3")

	_, ss, _ := c.Compile().GetResponse()
	r := strings.Join(ss, " ")
	if r != "before 1 AB middle 2 AC0E last 3" {
		t.Errorf("1) compile %s expected %s", r, "before 1 AB middle 2 AC0E last 3")
	}
}

func TestMask(t *testing.T) {
	jsonData := `{
		"cpf": 12345678910
	}`

	var data map[string]any
	json.Unmarshal([]byte(jsonData), &data)

	c := cp.NewCompiler().
		SetData(data).
		SetTxt("data.cpf").
		NewMaskAddOn("CPF")

	s, _, _ := c.Compile().GetResponse()
	if s != "123.456.789-10" {
		t.Errorf("1) compile %s expected %s", s, "123.456.789-10")
	}
}

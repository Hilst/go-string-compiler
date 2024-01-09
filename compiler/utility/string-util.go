package utility

import (
	"errors"
	"strings"

	l "me.compiler/compiler/lexer"
)

func splitPath(tk l.Token) []string {
	ids := make([]string, 0)
	for _, id := range strings.Split(tk.Literal, ".") {
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func Stringfy(d any) (string, error) {
	if d == nil {
		return StringInvalidEvent, errors.New("data run result is nil")
	}
	return d.(string), nil
}

func ApplyMask(toMask any, mask string) (string, error) {
	switch typed := toMask.(type) {
	case float64:
		switch mask {
		case "CPF":
			return directMask(int(typed), "###.###.###-##"), nil
		}
	}
	return Stringfy(toMask)
}

func directMask(value int, pattern string) string {
	decimals := int(value)
	var digit int
	runes := make([]rune, len(pattern))
	last := len(runes) - 1
	for i := last; i >= 0; i-- {
		r := pattern[i]
		if r == 35 { // #
			digit = decimals % 10
			decimals = decimals / 10
			runes[i] = rune(digit) + 48 // ASCII
		} else {
			runes[i] = rune(r)
		}
	}
	s := string(runes)
	return s
}

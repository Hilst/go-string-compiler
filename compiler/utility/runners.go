package utility

import (
	"errors"
	"strconv"
	"strings"

	l "me.compiler/compiler/lexer"
	p "me.compiler/compiler/parser"
)

func runData(path []string, data map[string]any) (any, error) {
	var d any
	d = data
	for _, p := range path {
		switch t := d.(type) {
		case []any:
			if idx, err := strconv.ParseInt(p, 0, 0); err == nil {
				d = t[idx]
				continue
			} else {
				return nil, err
			}
		case map[string]any:
			d = t[p]
		default:
			return d, nil
		}
	}
	return d, nil
}

// 0: data.
// 1: ID
func RunArray(event p.ParserEvent, data map[string]any) ([]any, error) {
	pArr, err := runData(splitPath(event.Tokens[1]), data)
	if err != nil {
		return nil, err
	}
	arr, castOK := pArr.([]any)
	if !castOK {
		return nil, errors.New("could not cast as splice")
	}
	return arr, nil
}

// 0: data.
// 1: ID
// 2: [
// 3: ID
// 4: ===
// 5: ID
// 6: ]
// 7: ID
func RunArraySelector(data map[string]any, tks []l.Token) (any, error) {
	castOK := false
	arrayPath := splitPath(tks[1])
	d, err := runData(arrayPath, data)
	if err != nil {
		return StringInvalidEvent, err
	}
	arr, castOK := d.([]any)
	if !castOK {
		return StringInvalidEvent, errors.New("is not array")
	}
	var comp any
	var looked any
	compPath := splitPath(tks[3])
	var m map[string]any

searchLoop:
	for _, obj := range arr {
		m, castOK = obj.(map[string]any)
		if !castOK {
			return StringInvalidEvent, errors.New("invalid map cast")
		}
		comp, err = runData(compPath, m)
		switch t := comp.(type) {
		case string:
			if t == tks[5].Literal {
				looked = obj
				break searchLoop
			}
		case float64:
			if f, err := strconv.ParseFloat(tks[5].Literal, 64); err == nil && f == comp {
				looked = obj
				break searchLoop
			}
		default:
			return StringInvalidEvent, errors.New("invalid comparator type")
		}
	}
	if looked == nil {
		return StringInvalidEvent, errors.New("not found any object")
	}
	targetPath := splitPath(tks[7])
	m, castOK = looked.(map[string]any)
	if !castOK {
		return StringInvalidEvent, errors.New("invalid map cast")
	}
	if r, err := runData(targetPath, m); err == nil {
		return r, nil
	}
	return StringInvalidEvent, err
}

// 0: data.
// 1: ID
func RunFreeCompile(data map[string]any, tks []l.Token) (any, error) {
	ids := splitPath(tks[1])
	d, err := runData(ids, data)
	if err != nil {
		return StringInvalidEvent, err
	} else if d == nil {
		return StringInvalidEvent, errors.New("illegal path empty value")
	}
	return d, err
}

// 0: [
// 1: count
// 2: |
// 3: singular
// 4: |
// 5: plural
// 6: ]
func RunCountLabel(arr []any, tks []l.Token) (string, error) {
	number := len(arr)
	if number <= 0 {
		return StringInvalidEvent, errors.New("invalid size array")
	}

	var name string
	if number == 1 {
		name = tks[3].Literal
	} else {
		name = tks[5].Literal
	}
	return strings.Join([]string{strconv.Itoa(number), name}, " "), nil
}

// 0: data.
// 1: ++
// 2: ID
func RunIterator(arr []any, tks []l.Token) ([]any, error) {
	values := make([]any, len(arr))
	var m map[string]any
	var castOk bool
	for i := 0; i < len(values); i++ {
		m, castOk = arr[i].(map[string]any)
		if !castOk {
			return values, errors.New("cant cast object to map")
		}
		v, err := runData(splitPath(tks[2]), m)
		if err != nil {
			return values, err
		}
		values[i], err = v, nil
		if err != nil {
			return values, err
		}
	}
	return values, nil
}

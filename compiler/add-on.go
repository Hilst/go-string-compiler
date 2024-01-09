package compiler

import (
	"errors"

	p "me.compiler/compiler/parser"
	u "me.compiler/compiler/utility"
)

type AddOnType int

const (
	ReadyAuxiliarArray AddOnType = iota
	ApplyMask
)

type PreAction func(ao AddOn, c *Compiler) (ok bool)
type PostAction func(ao AddOn, items []any) (results []any, err error)
type AddOn struct {
	type_          AddOnType
	masks          []string
	auxArrayString *string
	PreAction      PreAction
	PostAction     PostAction
}

func (c *Compiler) NewSetAuxArrayAddOn(p string) *Compiler {
	c.addOns.pre.addOns = append(c.addOns.pre.addOns,
		AddOn{
			type_:          ReadyAuxiliarArray,
			auxArrayString: &p,
			PreAction:      setAuxArray,
		},
	)

	return c
}

func (c *Compiler) NewMaskAddOn(masks ...string) *Compiler {
	c.addOns.post.addOns = append(c.addOns.post.addOns,
		AddOn{
			type_:      ApplyMask,
			masks:      masks,
			PostAction: applyMask,
		})
	return c
}

func setAuxArray(ao AddOn, c *Compiler) (ok bool) {
	if ao.type_ != ReadyAuxiliarArray {
		return false
	}
	c.input.psr.Start(*ao.auxArrayString)
	event := c.input.psr.Next()
	if event.EventType != p.MatchEvent || event.Match != p.FreeCompile {
		return false
	}
	arr, err := u.RunArray(event, c.input.data)
	if err != nil {
		return false
	}

	c.addOns.pre.auxiliarArray = arr
	return true
}

func applyMask(ao AddOn, items []any) (results []any, err error) {
	masked := make([]any, len(items))
	var mask string
	if len(ao.masks) != len(items) {
		return masked, errors.New("invalid mask slice size")
	}
	for i := 0; i < len(items); i++ {
		mask = ao.masks[i]
		masked[i], err = u.ApplyMask(items[i], mask)
		if err != nil {
			return masked, err
		}
	}
	return masked, nil
}

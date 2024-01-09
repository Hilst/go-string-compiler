package compiler

import (
	"errors"
	"strings"

	l "me.compiler/compiler/lexer"
	p "me.compiler/compiler/parser"
	u "me.compiler/compiler/utility"
)

type CompilerValue int

const (
	Data CompilerValue = iota
	PreAddOns
	PostAddOns
	Parser
)

type Compiler struct {
	input struct {
		psr    p.Parser
		orgTxt string
		data   map[string]any
	}
	addOns struct {
		pre struct {
			addOns        []AddOn
			auxiliarArray []any
		}
		post struct {
			addOns []AddOn
		}
	}
	output struct {
		parts        []any
		multiparts   bool
		processError error
	}
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Reset(cvs ...CompilerValue) *Compiler {
	if len(cvs) <= 0 {
		c.input.data = nil
		c.input.psr.Start(c.input.orgTxt)
		c.addOns.pre.auxiliarArray = make([]any, 0)
		c.addOns.pre.addOns = make([]AddOn, 0)
		c.addOns.post.addOns = make([]AddOn, 0)
		return c
	}
	for _, cv := range cvs {
		switch cv {
		case Data:
			c.input.data = nil
		case PreAddOns:
			c.addOns.pre.addOns = make([]AddOn, 0)
			c.addOns.pre.auxiliarArray = make([]any, 0)
		case PostAddOns:
			c.addOns.post.addOns = make([]AddOn, 0)
		case Parser:
			c.input.psr.Start(c.input.orgTxt)
		}
	}
	return c
}

func (c *Compiler) RunPreAddOns() *Compiler {
	for _, a := range c.addOns.pre.addOns {
		if !a.PreAction(a, c) {
			c.setResult(nil, false, errors.New("addon run invalid"))
		}
	}
	return c
}

func (c *Compiler) SetData(data map[string]any) *Compiler {
	c.input.data = data
	return c
}

func (c *Compiler) SetTxt(txt string) *Compiler {
	c.input.orgTxt = txt
	c.input.psr.Start(txt)
	return c
}

func (c *Compiler) Compile() *Compiler {
	if c.output.processError != nil {
		return c
	}

	result := make([]any, 0)

	for e := c.input.psr.Next(); e.EventType != p.EofEvent; e = c.input.psr.Next() {
		if e.EventType == p.InvalidEvent {
			result = append(result, u.StringInvalidEvent)
			c.setResult(result, false, errors.New(u.StringInvalidEvent))
			return c
		}

		switch e.EventType {
		case p.CommonEvent:
			common := make([]string, len(e.Tokens))
			for i := 0; i < len(common); i++ {
				common[i] = e.Tokens[i].Literal
			}
			result = append(result, strings.Join(common, " "))
		case p.MatchEvent:
			s, multiResults, calcErr := c.calculateMatch(e.Match, e.Tokens)
			if calcErr != nil {
				result = append(result, u.StringInvalidEvent)
				c.setResult(result, multiResults, calcErr)
				return c
			} else if multiResults {
				c.setResult(s, multiResults, nil)
				return c
			}
			result = append(result, s...)
		}
	}
	c.setResult(result, false, nil)
	return c
}

func (c *Compiler) calculateMatch(type_ p.MatchType, tokens []l.Token) ([]any, bool, error) {
	calculated := make([]any, 1)
	multiResults := false
	var e error
	switch type_ {
	case p.FreeCompile:
		calculated[0], e = u.RunFreeCompile(c.input.data, tokens)
	case p.ArraySelector:
		calculated[0], e = u.RunArraySelector(c.input.data, tokens)
	case p.Iterator:
		calculated, e = u.RunIterator(c.addOns.pre.auxiliarArray, tokens)
		multiResults = true
	case p.CountLabel:
		calculated[0], e = u.RunCountLabel(c.addOns.pre.auxiliarArray, tokens)
	default:
		calculated[0] = u.StringInvalidEvent
		e = errors.New("not supported match type")
	}
	return calculated, multiResults, e
}

func (c *Compiler) setResult(ps []any, mp bool, e error) {
	if e == nil {
		ps, e = c.runPosAddOns(ps)
	}
	c.output.parts = ps
	c.output.multiparts = mp
	c.output.processError = e
}

func (c *Compiler) runPosAddOns(ps []any) ([]any, error) {
	var err error
	for _, ao := range c.addOns.post.addOns {
		ps, err = ao.PostAction(ao, ps)
	}
	return ps, err
}

func (c *Compiler) IsOk() (bool, error) {
	if c.output.processError != nil {
		return false, c.output.processError
	}
	return true, nil
}

func (c *Compiler) GetResponse() (firstPart string, allParts []string, isMultipart bool) {
	strs := make([]string, len(c.output.parts))
	for i := 0; i < len(strs); i++ {
		if str, err := u.Stringfy(c.output.parts[i]); err == nil {
			strs[i] = str
		}
	}
	return strs[0], strs, c.output.multiparts
}

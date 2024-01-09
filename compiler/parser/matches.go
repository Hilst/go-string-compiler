package parser

import (
	l "me.compiler/compiler/lexer"
)

type MatchType string

const (
	FreeCompile   MatchType = "freeCompile"
	ArraySelector MatchType = "arraySelector"
	Iterator      MatchType = "iterator"
	CountLabel    MatchType = "countLabel"
	InvalidMatch  MatchType = "INVALID"
)

type matchObjective struct {
	tktypes_ []l.TokenType
	type_    MatchType
}

var allMatches = []matchObjective{
	{
		type_:    FreeCompile,
		tktypes_: []l.TokenType{l.Data, l.Ident},
	},
	{
		type_: ArraySelector,
		tktypes_: []l.TokenType{
			l.Data, l.Ident, l.LBracket,
			l.Ident, l.Congruent, l.Ident,
			l.RBracket, l.Ident},
	},
	{
		type_:    Iterator,
		tktypes_: []l.TokenType{l.Data, l.Iterator, l.Ident},
	},
	{
		type_: CountLabel,
		tktypes_: []l.TokenType{
			l.LBracket, l.Count, l.Pipeline,
			l.Ident, l.Pipeline,
			l.Ident, l.RBracket,
		},
	},
}

func (ps *Parser) findMatchEvent() ParserEvent {
	avaible := make([]matchObjective, 0)
	for _, obj := range allMatches {
		if len(obj.tktypes_) == len(ps.currTks) {
			avaible = append(avaible, obj)
		}
	}

	var match *MatchType
	index := 0
	for index < len(ps.currTks) && len(avaible) > 0 {
		for _, poss := range avaible {
			if poss.tktypes_[index] == ps.currTks[index].Type {
				if index == len(poss.tktypes_)-1 {
					match = &poss.type_
				}
			}
		}
		index++
	}

	if match == nil {
		return CreateEvent(InvalidEvent, InvalidMatch, ps.currTks)
	}
	return CreateEvent(MatchEvent, *match, ps.currTks)
}

package parser

import (
	l "me.compiler/compiler/lexer"
)

type ParserEventType string

const (
	MatchEvent   ParserEventType = "match"
	CommonEvent  ParserEventType = "common"
	EofEvent     ParserEventType = "eof"
	InvalidEvent ParserEventType = "invalid"
)

type ParserEvent struct {
	EventType ParserEventType
	Match     MatchType
	Tokens    []l.Token
}

type Parser struct {
	lx               l.Lexer
	firstTk          *l.TokenType
	currTks          []l.Token
	currentNeedMatch bool
}

func NewParser() Parser {
	return Parser{
		lx:               l.NewLexer(),
		firstTk:          nil,
		currTks:          make([]l.Token, 0),
		currentNeedMatch: false,
	}
}

func CreateEvent(
	EventType ParserEventType,
	Match MatchType,
	Tokens []l.Token,
) ParserEvent {
	return ParserEvent{EventType, Match, Tokens}
}

func (ps *Parser) Start(s string) {
	ps.lx.Start(s)
	ps.firstTk = nil
	ps.currTks = make([]l.Token, 0)
	ps.currentNeedMatch = false
}

func (ps *Parser) Next() ParserEvent {
	t := ps.lx.Next()
	if t.Type == l.Eof {
		return CreateEvent(EofEvent, InvalidMatch, make([]l.Token, 0))
	}

	if ps.firstTk == nil {
		ps.firstTk = &t.Type
	}
	ps.currTks = append(ps.currTks, t)

	if *ps.firstTk == l.Data || *ps.firstTk == l.LBracket {
		return ps.parseDataOnly()
	}
	return ps.parseWithConcat()
}

func (ps *Parser) parseWithConcat() ParserEvent {
	defer ps.resetCurrentTokens()
	foundConcat, isConcatOpen := ps.saveNextTokens()
	ps.updateConcatStatus(foundConcat, isConcatOpen)

	if ps.currentNeedMatch {
		return ps.findMatchEvent()
	}
	return CreateEvent(CommonEvent, InvalidMatch, ps.currTks)
}

func (ps *Parser) resetCurrentTokens() {
	ps.currTks = make([]l.Token, 0)
}

func (ps *Parser) updateConcatStatus(found bool, open bool) {
	if !found {
		if ps.currentNeedMatch {
			ps.currentNeedMatch = false
		}
		return
	}
	if !ps.currentNeedMatch && !open {
		ps.currentNeedMatch = true
		return
	}
	if ps.currentNeedMatch && open {
		ps.currentNeedMatch = false
		return
	}
}

func (ps *Parser) parseDataOnly() ParserEvent {
	defer ps.resetCurrentTokens()
	ps.saveNextTokens()
	return ps.findMatchEvent()
}

func (ps *Parser) saveNextTokens() (foundConcat bool, isConcatOpen bool) {
	tk := ps.lx.Next()
	for !tokenOfTypes(tk.Type, l.Illegal, l.Eof) {
		ps.currTks = append(ps.currTks, tk)
		tk = ps.lx.Next()
		if tk.Type == l.LConcat {
			return true, true
		} else if tk.Type == l.RConcat {
			return true, false
		}
	}
	return false, false
}

func tokenOfTypes(tt l.TokenType, types ...l.TokenType) bool {
	for _, type_ := range types {
		if type_ == tt {
			return true
		}
	}
	return false
}

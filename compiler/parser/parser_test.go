package parser

import (
	"testing"

	l "me.compiler/compiler/lexer"
)

func TestFreeCall(t *testing.T) {
	const input = "data.ident"

	var sut = NewParser()
	sut.Start(input)

	var expected []ParserEvent
	expected = append(expected,
		CreateEvent(MatchEvent, FreeCompile, []l.Token{
			l.CreateToken(l.Data, "data."),
			l.CreateToken(l.Ident, "ident"),
		}),
	)

	var out ParserEvent
	var eqTks bool
	var tryTk l.Token
	for _, exp := range expected {
		out = sut.Next()
		eqTks = true
		for indx, tk := range out.Tokens {
			if indx < len(exp.Tokens) {
				tryTk = exp.Tokens[indx]
				if tk != tryTk {
					eqTks = false
					break
				}
			} else {
				eqTks = false
				break
			}
		}

		if out.EventType != exp.EventType || out.Match != exp.Match || !eqTks {
			t.Errorf("P Event %s expected %s", out, exp)
		}
	}
}

func TestIterator(t *testing.T) {
	const input = "data.++.ident"

	var sut = NewParser()
	sut.Start(input)

	var expected []ParserEvent
	expected = append(expected,
		CreateEvent(MatchEvent, Iterator, []l.Token{
			l.CreateToken(l.Data, "data."),
			l.CreateToken(l.Iterator, "++"),
			l.CreateToken(l.Ident, ".ident"),
		}),
	)

	var out ParserEvent
	var eqTks bool
	var tryTk l.Token
	for _, exp := range expected {
		out = sut.Next()
		eqTks = true
		for indx, tk := range out.Tokens {
			if indx < len(exp.Tokens) {
				tryTk = exp.Tokens[indx]
				if tk != tryTk {
					eqTks = false
					break
				}
			} else {
				eqTks = false
				break
			}
		}

		if out.EventType != exp.EventType || out.Match != exp.Match || !eqTks {
			t.Errorf("P Event %s expected %s", out, exp)
		}
	}
}

func TestArraySelector(t *testing.T) {
	const input = "data.identA.identB[identC===identD].identE"

	var sut = NewParser()
	sut.Start(input)

	var tokens []l.Token
	tokens = append(tokens,
		l.CreateToken(l.Data, "data."),
		l.CreateToken(l.Ident, "identA.identB"),
		l.CreateToken(l.LBracket, "["),
		l.CreateToken(l.Ident, "identC"),
		l.CreateToken(l.Congruent, "==="),
		l.CreateToken(l.Ident, "identD"),
		l.CreateToken(l.RBracket, "]"),
		l.CreateToken(l.Ident, ".identE"),
	)

	var expected []ParserEvent
	expected = append(expected,
		CreateEvent(MatchEvent, ArraySelector, tokens),
	)

	var out ParserEvent
	var eqTks bool
	var tryTk l.Token
	for _, exp := range expected {
		out = sut.Next()
		eqTks = true
		for indx, tk := range out.Tokens {
			if indx < len(exp.Tokens) {
				tryTk = exp.Tokens[indx]
				if tk != tryTk {
					eqTks = false
					break
				}
			} else {
				eqTks = false
				break
			}
		}

		if out.EventType != exp.EventType || out.Match != exp.Match || !eqTks {
			t.Errorf("P Event %s expected %s", out, exp)
		}
	}
}

func TestCountExpression(t *testing.T) {
	const input = "[count|singular|plural]"

	var tokens []l.Token
	tokens = append(tokens,
		l.CreateToken(l.LBracket, "["),
		l.CreateToken(l.Count, "count"),
		l.CreateToken(l.Pipeline, "|"),
		l.CreateToken(l.Ident, "singular"),
		l.CreateToken(l.Pipeline, "|"),
		l.CreateToken(l.Ident, "plural"),
		l.CreateToken(l.RBracket, "]"),
	)

	var sut = NewParser()
	sut.Start(input)

	var expected []ParserEvent
	expected = append(expected,
		CreateEvent(MatchEvent, CountLabel, tokens),
	)

	var out ParserEvent
	var eqTks bool
	var tryTk l.Token
	for _, exp := range expected {
		out = sut.Next()
		eqTks = true
		for indx, tk := range out.Tokens {
			if indx < len(exp.Tokens) {
				tryTk = exp.Tokens[indx]
				if tk != tryTk {
					eqTks = false
					break
				}
			} else {
				eqTks = false
				break
			}
		}

		if out.EventType != exp.EventType || out.Match != exp.Match || !eqTks {
			t.Errorf("P Event %s expected %s", out, exp)
		}
	}
}

func TestWithConcat(t *testing.T) {
	const input = "first second @{data.ident0} one middle @{data.++.ident1} two middle @{data.identA.identB[identC===identD].identE} before last"

	var sut = NewParser()
	sut.Start(input)

	var asTks []l.Token
	asTks = append(asTks,
		l.CreateToken(l.Data, "data."),
		l.CreateToken(l.Ident, "identA.identB"),
		l.CreateToken(l.LBracket, "["),
		l.CreateToken(l.Ident, "identC"),
		l.CreateToken(l.Congruent, "==="),
		l.CreateToken(l.Ident, "identD"),
		l.CreateToken(l.RBracket, "]"),
		l.CreateToken(l.Ident, ".identE"),
	)

	var expected []ParserEvent
	expected = append(expected,
		CreateEvent(CommonEvent, InvalidMatch, []l.Token{
			l.CreateToken(l.Ident, "first"),
			l.CreateToken(l.Ident, "second"),
		}),
		CreateEvent(MatchEvent, FreeCompile, []l.Token{
			l.CreateToken(l.Data, "data."),
			l.CreateToken(l.Ident, "ident0"),
		}),
		CreateEvent(CommonEvent, InvalidMatch, []l.Token{
			l.CreateToken(l.Ident, "one"),
			l.CreateToken(l.Ident, "middle"),
		}),
		CreateEvent(MatchEvent, Iterator, []l.Token{
			l.CreateToken(l.Data, "data."),
			l.CreateToken(l.Iterator, "++"),
			l.CreateToken(l.Ident, ".ident1"),
		}),
		CreateEvent(CommonEvent, InvalidMatch, []l.Token{
			l.CreateToken(l.Ident, "two"),
			l.CreateToken(l.Ident, "middle"),
		}),
		CreateEvent(MatchEvent, ArraySelector, asTks),
		CreateEvent(CommonEvent, InvalidMatch, []l.Token{
			l.CreateToken(l.Ident, "before"),
			l.CreateToken(l.Ident, "last"),
		}),
	)

	var out ParserEvent
	var eqTks bool
	var tryTk l.Token
	for _, exp := range expected {
		out = sut.Next()
		eqTks = true
		for indx, tk := range out.Tokens {
			if indx < len(exp.Tokens) {
				tryTk = exp.Tokens[indx]
				if tk != tryTk {
					eqTks = false
					break
				}
			} else {
				eqTks = false
				break
			}
		}

		if out.EventType != exp.EventType || out.Match != exp.Match || !eqTks {
			t.Errorf("P Event %s expected %s", out, exp)
		}
	}
}

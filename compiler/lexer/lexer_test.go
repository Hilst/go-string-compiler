package lexer

import "testing"

func TestBase(t *testing.T) {
	const input = "lorem ipsum"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(Ident, "lorem"),
		CreateToken(Ident, "ipsum"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

func TestFreeCall(t *testing.T) {
	const input = "data.ident"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(Data, "data."),
		CreateToken(Ident, "ident"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

func TestIterator(t *testing.T) {
	const input = "data.++.ident"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(Data, "data."),
		CreateToken(Iterator, "++"),
		CreateToken(Ident, ".ident"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

func TestArraySelector(t *testing.T) {
	const input = "data.identA.identB[identC===identD].identE"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(Data, "data."),
		CreateToken(Ident, "identA.identB"),
		CreateToken(LBracket, "["),
		CreateToken(Ident, "identC"),
		CreateToken(Congruent, "==="),
		CreateToken(Ident, "identD"),
		CreateToken(RBracket, "]"),
		CreateToken(Ident, ".identE"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

func TestCountExpression(t *testing.T) {
	const input = "[count|singular|plural]"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(LBracket, "["),
		CreateToken(Count, "count"),
		CreateToken(Pipeline, "|"),
		CreateToken(Ident, "singular"),
		CreateToken(Pipeline, "|"),
		CreateToken(Ident, "plural"),
		CreateToken(RBracket, "]"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

func TestAllAndConcat(t *testing.T) {
	const input = "first @{data.ident0} one @{data.++.ident1} two @{data.identA.identB[identC===identD].identE} last"

	var sut = NewLexer()
	sut.Start(input)

	var expected []Token
	expected = append(expected,
		CreateToken(Ident, "first"),
		CreateToken(LConcat, "@{"),
		CreateToken(Data, "data."),
		CreateToken(Ident, "ident0"),
		CreateToken(RConcat, "}"),
		CreateToken(Ident, "one"),
		CreateToken(LConcat, "@{"),
		CreateToken(Data, "data."),
		CreateToken(Iterator, "++"),
		CreateToken(Ident, ".ident1"),
		CreateToken(RConcat, "}"),
		CreateToken(Ident, "two"),
		CreateToken(LConcat, "@{"),
		CreateToken(Data, "data."),
		CreateToken(Ident, "identA.identB"),
		CreateToken(LBracket, "["),
		CreateToken(Ident, "identC"),
		CreateToken(Congruent, "==="),
		CreateToken(Ident, "identD"),
		CreateToken(RBracket, "]"),
		CreateToken(Ident, ".identE"),
		CreateToken(RConcat, "}"),
		CreateToken(Ident, "last"),
	)

	var out Token
	for _, exp := range expected {
		out = sut.Next()
		if out != exp {
			t.Errorf("Token %s expected %s", out, exp)
		}
	}
}

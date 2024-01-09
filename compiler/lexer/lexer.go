package lexer

type TokenType string

const (
	Illegal TokenType = "ILLEGAL"
	Eof     TokenType = "EOF"
	Ident   TokenType = "IDENT"

	LConcat TokenType = "@{"
	RConcat TokenType = "}"

	LBracket TokenType = "["
	RBracket TokenType = "]"

	Congruent TokenType = "==="
	Pipeline  TokenType = "|"
	Iterator  TokenType = "++"

	Data  TokenType = "data."
	Count TokenType = "count"
)

const (
	_at       = int('@')
	_lsquigly = int('{')
	_rsquigly = int('}')

	_lbracket = int('[')
	_rbracket = int(']')

	_equal = int('=')
	_pipe  = int('|')
	_plus  = int('+')

	_whiteSpace = int(' ')
	_eof        = int('\x00')
)

type Token struct {
	Type    TokenType
	Literal string
}

func CreateToken(type_ TokenType, literal string) Token {
	return Token{type_, literal}
}

type Lexer struct {
	readPos int
	ch      rune
	s       string
}

func NewLexer() Lexer {
	return Lexer{}
}

func isValidIdent(character rune) bool {
	char := int(character)
	return !(char == _at ||
		char == _equal ||
		char == _lbracket ||
		char == _rbracket ||
		char == _lsquigly ||
		char == _rsquigly ||
		char == _pipe ||
		char == _plus ||
		char == _whiteSpace ||
		char == _eof)
}

var keyWords = map[string]Token{
	"data.": CreateToken(Data, "data."),
	"count": CreateToken(Count, "count"),
}

func (lx *Lexer) Start(s string) {
	lx.s = s
	lx.readPos = 0
	lx.readChar()
}

func (lx *Lexer) Next() Token {
	lx.skipWhitespace()

	var tk Token
	var tkNil bool = true
	switch lx.ch {
	case '@':
		if lx.peekChar(0) == rune(_lsquigly) {
			lx.readChar()
			tk = CreateToken(LConcat, "@{")
			tkNil = false
		} else {
			tk = CreateToken(Illegal, "@")
			tkNil = false
		}
	case '}':
		tk = CreateToken(RConcat, "}")
		tkNil = false
	case '[':
		tk = CreateToken(LBracket, "[")
		tkNil = false
	case ']':
		tk = CreateToken(RBracket, "]")
		tkNil = false
	case '=':
		if lx.peekChar(0) == rune(_equal) && lx.peekChar(1) == rune(_equal) {
			lx.readChar()
			lx.readChar()
			tk = CreateToken(Congruent, "===")
			tkNil = false
		}
	case '+':
		if lx.peekChar(0) == rune(_plus) {
			lx.readChar()
			tk = CreateToken(Iterator, "++")
			tkNil = false
		}
	case '|':
		tk = CreateToken(Pipeline, "|")
		tkNil = false
	case '\x00':
		tk = CreateToken(Eof, "eof")
		tkNil = false
	}

	if isValidIdent(lx.ch) {
		ident := lx.readIdent()
		if kw, ok := keyWords[ident]; ok {
			return kw
		} else {
			return CreateToken(Ident, ident)
		}
	} else if tkNil {
		return CreateToken(Illegal, string(lx.ch))
	}

	lx.readChar()
	return tk
}

func (lx *Lexer) readChar() {
	if lx.readPos >= len(lx.s) {
		lx.ch = '\x00'
	} else {
		lx.ch = rune(lx.s[lx.readPos])
	}
	lx.readPos++
}

func (lx *Lexer) skipWhitespace() {
	for lx.ch == ' ' {
		lx.readChar()
	}
}

func (lx *Lexer) peekChar(upfront int) rune {
	if lx.readPos+upfront >= len(lx.s) {
		return '\x00'
	} else {
		return rune(lx.s[lx.readPos+upfront])
	}
}

func (lx *Lexer) readIdent() string {
	curr := ""
	for isValidIdent(lx.ch) {
		curr += string(lx.ch)
		_, isKw := keyWords[curr]
		lx.readChar()
		if isKw {
			return curr
		}
	}
	return curr
}

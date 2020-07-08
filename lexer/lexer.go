package lexer

import (
	"simpsel/token"
	"strings"
)

type Lexer struct {
	input        string // The input to be lexed
	position     int    // Current position in input (current char)
	readPosition int    // Current reading position in input (after current char)
	ch           byte   // Current byte under examination
	line         int    // The line number
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '#':
		l.readChar()
		if num := l.readNumber(); num != "" {
			tok.Type = token.INT
			tok.Literal = num
			tok.Line = l.line
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	case '$':
		l.readChar()
		if num := l.readNumber(); num != "" {
			tok.Type = token.REGISTER
			tok.Literal = num
			tok.Line = l.line
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	case '.':
		l.readChar()
		if isVarTer(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupDirective(strings.ToLower(tok.Literal))
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	case ';':
		tok = newToken(token.COMMENT, l.ch, l.line)
		l.skipUntilNewline()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = l.line
	default:
		if isVarTer(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(strings.ToLower(tok.Literal))
			tok.Line = l.line
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.line)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	if l.ch == '\n' {
		l.line ++
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isVarTer(l.ch) {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipUntilNewline() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func isVarTer(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tType token.TokenType, ch byte, line int) token.Token {
	return token.Token{Type: tType, Literal: string(ch), Line: line}
}
package main

import (
	"io/ioutil"
	"os"
)

func main() {
	input, error := ioutil.ReadFile(os.Args[1])
	if error != nil {
		println("Error reading file")
		os.Exit(1)
	}
	l := Lexer{input: string(input)}
	l.readChar()
	for {
		tok := l.NextToken()
		if tok.Type == ILLEGAL {
			break
		}
		println(tok.Type, tok.Literal)
	}
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	IDENT   = "IDENT"
	COMMA   = ","
	LPAREN  = "("
	RPAREN  = ")"
	INT     = "INT"
	WEEKDAY = "WEEKDAY"
	UNIT    = "UNIT"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

var keywords = map[string]TokenType{
        "monday": WEEKDAY,
        "tuesday": WEEKDAY,
        "wednesday": WEEKDAY,
        "thursday": WEEKDAY,
        "friday": WEEKDAY,
        "saturday": WEEKDAY,
        "sunday": WEEKDAY,
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	switch l.ch {
	case ',':
		tok = newToken(COMMA, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			if t, ok := keywords[tok.Literal]; ok {
                                tok.Type = t
			} else {
                                tok.Type = IDENT
                        }

			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
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

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

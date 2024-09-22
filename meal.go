package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	input, error := os.ReadFile(os.Args[1])
	if error != nil {
		println("Error reading file")
		os.Exit(1)
	}
	l := Lexer{input: string(input)}
	p := Parser{l: &l}
	program := p.Parse()

        run(program)
}

func run(program *Program) {

        meals := make(map[string]int)
        recipes := make(map[string]map[string]int)
        units := make(map[string]string)
        result := make(map[string]int)

	for _, statement := range program.Statements {
		switch s := statement.(type) {
		case *WeekdayStatement:
			for _, s := range s.Meals {

				switch s := s.(type) {
				case *MealStatement:
					if _, ok := meals[s.Name]; !ok {
						meals[s.Name] = 0
					}
					meals[s.Name] += 1
                                default:
                                        fmt.Printf("Error: unexpected statement in weekday. Expected MealStatement, found %T\n", s)
                                        os.Exit(1)
				}
			}
		case *RecipeStatement:
			for _, i := range s.Ingredients {
				switch i := i.(type) {
				case *IngredientStatement:
					if _, ok := recipes[s.Name]; !ok {
						recipes[s.Name] = make(map[string]int)
					}
					recipes[s.Name][i.Name] += i.Amount

                                        if _, ok := units[i.Name]; !ok {
                                                units[i.Name] = i.Unit
                                        }


                                default:
                                        fmt.Printf("Error: unexpected statement in weekday. Expected IngredientStatement, found %T\n", s)
                                        os.Exit(1)
				}

			}
		}
	}

        for meal, count := range meals {
                recipe := recipes[meal]
                for ingredient, amount := range recipe {
                        if _, ok := result[ingredient]; !ok {
                                result[ingredient] = 0
                        }
                        result[ingredient] += amount * count
                }

        }

        fmt.Println("Shopping list:")
        for ingredient, amount := range result {
                fmt.Printf("- %s %d%s\n", ingredient, amount, units[ingredient])
        }

}

type Node interface {
	TokenLiteral() string
	String() string
}
type Statement interface {
	Node
	statementNode()
}
type Program struct {
	Statements []Statement
}
type WeekdayStatement struct {
	token   Token
	Weekday string
	Meals   []Statement
}
type MealStatement struct {
	token Token
	Name  string
}

func (ms *MealStatement) statementNode()       {}
func (ms *MealStatement) TokenLiteral() string { return ms.token.Literal }
func (ms *MealStatement) String() string {
	return ms.token.Literal
}
func (ws *WeekdayStatement) statementNode()       {}
func (ws *WeekdayStatement) TokenLiteral() string { return ws.token.Literal }
func (ws *WeekdayStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ws.token.Literal)
	out.WriteString(" ")
	for _, s := range ws.Meals {
		out.WriteString(s.String())
	}
	return out.String()
}

type RecipeStatement struct {
	token       Token
	Name        string
	Ingredients []Statement
}

func (rs *RecipeStatement) statementNode()       {}
func (rs *RecipeStatement) TokenLiteral() string { return rs.token.Literal }
func (rs *RecipeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.token.Literal)
	out.WriteString(" ")
	out.WriteString(rs.Name)
	out.WriteString(" ")
	for _, s := range rs.Ingredients {
		out.WriteString(s.String())
		out.WriteString(" ")
	}
	return out.String()
}

type IngredientStatement struct {
	token  Token
	Name   string
	Amount int
	Unit   string
}

func (is *IngredientStatement) statementNode()       {}
func (is *IngredientStatement) TokenLiteral() string { return is.token.Literal }
func (is *IngredientStatement) String() string {
	var out bytes.Buffer
	out.WriteString(is.token.Literal)
	out.WriteString(" ")
	out.WriteString(strconv.Itoa(is.Amount))
	out.WriteString(is.Unit)
	return out.String()
}

type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Parser) Parse() *Program {
	program := &Program{}
	program.Statements = []Statement{}
	p.l.readChar()
	p.nextToken()
	p.nextToken()
	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.nextToken()
}
func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case WEEKDAY:
		return p.parseWeekdayStatement()
	case IDENT:
		if p.peekTokenIs(LPAREN) {
			return p.parseRecipeStatement()
		} else {
			if p.peekTokenIs(INT) {
				return p.parseIngredientStatement()
			} else {
				return p.parseMealStatement()
			}
		}
	default:
		return nil
	}
}

func (p *Parser) parseMealStatement() *MealStatement {
	stmt := &MealStatement{token: p.curToken}
	stmt.Name = p.curToken.Literal
	return stmt
}

func (p *Parser) parseWeekdayStatement() *WeekdayStatement {
	stmt := &WeekdayStatement{token: p.curToken}
	if !p.expectPeek(LPAREN) {
		return nil
	}
	for !p.peekTokenIs(RPAREN) {
		p.nextToken()
		stmt.Meals = append(stmt.Meals, p.parseStatement())
	}
	p.nextToken()
	return stmt
}
func (p *Parser) parseRecipeStatement() *RecipeStatement {
	stmt := &RecipeStatement{token: p.curToken}
	stmt.Name = p.curToken.Literal
	if !p.expectPeek(LPAREN) {
		return nil
	}
	for !p.peekTokenIs(RPAREN) {
		p.nextToken()
		stmt.Ingredients = append(stmt.Ingredients, p.parseStatement())
	}
	p.nextToken()
	return stmt
}

func (p *Parser) parseIngredientStatement() *IngredientStatement {
	stmt := &IngredientStatement{token: p.curToken}
	stmt.Name = p.curToken.Literal
	p.nextToken()
	i, error := strconv.Atoi(p.curToken.Literal)
	if error != nil {
		println("Error parsing amount for ingredient " + stmt.Name)
		fmt.Printf("Error: %v\n", error)
		os.Exit(1)
	}
	stmt.Amount = i
	p.nextToken()
	stmt.Unit = p.curToken.Literal
	return stmt
}

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
	INT     = "INT"
	LPAREN  = "("
	RPAREN  = ")"
	UNIT    = "UNIT"
	WEEKDAY = "WEEKDAY"
	IDENT   = "IDENT"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

var keywords = map[string]TokenType{
	"monday":    WEEKDAY,
	"tuesday":   WEEKDAY,
	"wednesday": WEEKDAY,
	"thursday":  WEEKDAY,
	"friday":    WEEKDAY,
	"saturday":  WEEKDAY,
	"sunday":    WEEKDAY,
}
var units = map[string]TokenType{
	"u":    UNIT,
	"g":    UNIT,
	"kg":   UNIT,
	"ml":   UNIT,
	"l":    UNIT,
	"cup":  UNIT,
	"tbsp": UNIT,
}

func (l *Lexer) nextToken() Token {
	var tok Token
	l.skipWhitespace()
	switch l.ch {
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			if t, ok := keywords[tok.Literal]; ok {
				tok.Type = t
			} else if t, ok := units[tok.Literal]; ok {
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

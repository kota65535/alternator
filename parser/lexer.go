package parser

import (
	"fmt"
	"github.com/kota65535/alternator/lexer"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type Parser struct {
	lexer     *lexer.Lexer
	lastToken *lexer.Token
	result    []Statement
	lastError error
}

func NewParser(reader io.Reader) *Parser {

	tokens := []lexer.TokenType{}
	skippedTokens := []lexer.TokenType{}

	for k, v := range Keywords {
		tokens = append(tokens, lexer.NewSimpleTokenType(lexer.TokenID(k), v, true, 1))
	}
	for k, v := range Signs {
		tokens = append(tokens, lexer.NewSimpleTokenType(lexer.TokenID(k), v, true, 1))
	}
	for _, v := range Literals {
		tokens = append(tokens, v)
	}

	// Keywords that contain "NOT" token can lead to fail parse error in column options for example:
	//
	//   { [NOT] NULL | CHECK(...) [NOT] ENFORCED } ...
	//
	// And assume the following statement:
	//
	//   CHECK(...) NOT NULL
	//
	// In this case, shift takes precedence and ENFORCED token is expected, so NULL token causes parse error.
	// To prevent this, we treat keyword containing "NOT" (ex: "NOT NULL") as a single token, like "NOT\s+NULL".
	tokens = append(tokens, lexer.NewRegexpTokenType(NOT_ENFORCED, "(NOT|not)\\s+(ENFORCED|enforced)", 1))

	for _, v := range Skipped {
		skippedTokens = append(skippedTokens, lexer.NewRegexpTokenType(-1, v, 0))
	}

	l := lexer.NewLexer(reader, tokens, skippedTokens)

	return &Parser{
		lexer: l,
	}
}

func (p *Parser) Parse() ([]Statement, error) {
	ret := yyParse(p)
	if ret != 0 {
		return nil, fmt.Errorf("parse failed: %w", p.LastError())
	}
	return p.result, nil
}

func (p *Parser) Lex(lval *yySymType) int {
	token, err := p.lexer.Scan()
	if err != nil {
		p.Error(err.Error())
		return 0
	}
	if token == nil {
		return 0
	}

	lval.token = token

	p.lastToken = token

	logrus.Debugf("token '%s' as %s\n", token.Literal, token.Type.GetID())

	return int(token.Type.GetID())
}

func (p *Parser) Error(e string) {
	line, lineNum := p.lexer.GetLastLine()
	marks := strings.Repeat(" ", p.lastToken.Position.Column) + strings.Repeat("^", len(p.lastToken.Literal))
	p.lastError = fmt.Errorf("%s at line %d\n\n%s\n%s", e, lineNum, line, marks)
}

func (p *Parser) LastError() error {
	return p.lastError
}

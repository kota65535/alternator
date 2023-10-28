package parser

import (
	"errors"
	"fmt"
	"github.com/kota65535/alternator/lexer"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

type Parser struct {
	lexer     *lexer.Lexer
	lastToken *lexer.Token
	result    []Statement
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
		return nil, errors.New("parse failed")
	}
	return p.result, nil
}

func (p *Parser) Lex(lval *yySymType) int {
	token, err := p.lexer.Scan()
	if err != nil {
		if e, ok := err.(lexer.UnknownTokenError); ok {
			fmt.Fprintln(os.Stderr, e.Error()+":")
			fmt.Fprintln(os.Stderr, p.lexer.GetLastLine())
			fmt.Fprintln(os.Stderr, strings.Repeat(" ", e.Position.Column)+strings.Repeat("^", len(e.Literal)))
		}
		p.Error(err.Error())
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
	fmt.Fprintln(os.Stderr, e+":")
	fmt.Fprintln(os.Stderr, p.lexer.GetLastLine())
	fmt.Fprintln(os.Stderr, strings.Repeat(" ", p.lastToken.Position.Column)+strings.Repeat("^", len(p.lastToken.Literal)))
}

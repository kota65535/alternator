package lexer

import (
	"regexp"
	"strconv"
	"strings"
)

// TokenID is Identifier for TokenType.
type TokenID int

func (id TokenID) String() string {
	return strconv.Itoa(int(id))
}

type TokenType interface {
	GetID() TokenID
	FindToken(string, Position) *Token
	GetPriority() int
}

type SimpleTokenType struct {
	ID              TokenID
	Pattern         string
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewSimpleTokenType(id TokenID, pattern string, ignoreCase bool, priority int) *SimpleTokenType {
	return &SimpleTokenType{
		ID:         id,
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
		Priority:   priority,
	}
}

func (ptt *SimpleTokenType) String() string {
	return ptt.ID.String()
}

func (ptt *SimpleTokenType) GetID() TokenID {
	return ptt.ID
}

func (ptt *SimpleTokenType) FindToken(s string, p Position) *Token {
	if ptt.IgnoreCase {
		s = strings.ToUpper(s)
	}
	if !strings.HasPrefix(s, ptt.Pattern) {
		return nil
	}

	return &Token{
		Type:     ptt,
		Literal:  ptt.Pattern,
		Position: p,
	}
}

func (ptt *SimpleTokenType) GetPriority() int {
	return ptt.Priority
}

type RegexpTokenType struct {
	ID              TokenID
	Re              *regexp.Regexp
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewRegexpTokenType(id TokenID, re string, ignoreCase bool, priority int) *RegexpTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &RegexpTokenType{
		ID:         id,
		Re:         regexp.MustCompile(re),
		IgnoreCase: ignoreCase,
		Priority:   priority,
	}
}

func (rtt *RegexpTokenType) String() string {
	return rtt.ID.String()
}

func (rtt *RegexpTokenType) GetID() TokenID {
	return rtt.ID
}

func (rtt *RegexpTokenType) FindToken(s string, p Position) *Token {
	m := rtt.Re.FindStringSubmatch(s)
	if rtt.IgnoreCase {
		s = strings.ToUpper(s)
	}
	if len(m) > 0 {
		return &Token{
			Type:       rtt,
			Literal:    m[0],
			Submatches: m[1:],
			Position:   p,
		}
	}
	return nil
}

func (rtt *RegexpTokenType) GetPreviousTokenID() TokenID {
	return rtt.PreviousTokenID
}

func (rtt *RegexpTokenType) GetPriority() int {
	return rtt.Priority
}

type Token struct {
	Type       TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
	Position   Position // Position of token.
}

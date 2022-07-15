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

// BalancedParenthesesTokenType implementation

type BalancedParenthesesTokenType struct {
	ID          TokenID
	Re          *regexp.Regexp
	Parentheses []rune
	IgnoreCase  bool
	Priority    int
}

func NewBalancedParenthesesTokenType(id TokenID, re string, open rune, close rune, ignoreCase bool, priority int) *BalancedParenthesesTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &BalancedParenthesesTokenType{
		ID:          id,
		Re:          regexp.MustCompile(re),
		Parentheses: []rune{open, close},
		IgnoreCase:  ignoreCase,
		Priority:    priority,
	}
}

func (btt *BalancedParenthesesTokenType) String() string {
	return btt.ID.String()
}

func (btt *BalancedParenthesesTokenType) GetID() TokenID {
	return btt.ID
}

func (btt *BalancedParenthesesTokenType) FindToken(s string, p Position) *Token {
	if btt.IgnoreCase {
		s = strings.ToUpper(s)
	}

	if !strings.HasPrefix(s, string(btt.Parentheses[0])) {
		return nil
	}

	level := 0
	curIdx := 0
	for i, c := range s {
		parenthesesMatched := false
		if c == btt.Parentheses[0] {
			level += 1
			parenthesesMatched = true
		}
		if c == btt.Parentheses[1] {
			level += -1
			parenthesesMatched = true
		}
		if parenthesesMatched {
			m := btt.Re.FindStringSubmatch(s[curIdx:i])
			if len(m) == 0 {
				// Pattern not matched
				return nil
			}
			curIdx = i + 1
		}
		if level == 0 {
			break
		}
	}

	if level != 0 {
		// Parentheses not matched
		return nil
	}

	return &Token{
		Type:     btt,
		Literal:  s[:curIdx],
		Position: p,
	}
}

func (btt *BalancedParenthesesTokenType) GetPriority() int {
	return btt.Priority
}

type SyntaxAwareTokenType struct {
	startPattern []TokenID
	endPattern   []TokenID
	tokenType    TokenType
	Enabled      bool
}

func NewSyntaxAwareTokenType(startPattern []TokenID, endPattern []TokenID, tokenType TokenType) *SyntaxAwareTokenType {
	return &SyntaxAwareTokenType{
		startPattern: startPattern,
		endPattern:   endPattern,
		tokenType:    tokenType,
	}
}

func (btt *SyntaxAwareTokenType) GetPriority() int {
	return btt.tokenType.GetPriority()
}

type Token struct {
	Type       TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
	Position   Position // Position of token.
}

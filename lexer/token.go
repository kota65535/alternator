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
	GetPreviousTokenID() TokenID
	GetPriority() int
}

type SimpleTokenType struct {
	ID              TokenID
	Pattern         string
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewSimpleTokenType(id TokenID, pattern string, ignoreCase bool, priority int, previousTokenId TokenID) *SimpleTokenType {
	return &SimpleTokenType{
		ID:              id,
		Pattern:         pattern,
		IgnoreCase:      ignoreCase,
		Priority:        priority,
		PreviousTokenID: previousTokenId,
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

func (ptt *SimpleTokenType) GetPreviousTokenID() TokenID {
	return ptt.PreviousTokenID
}

func (ptt *SimpleTokenType) GetPriority() int {
	return ptt.Priority
}

type MultiTokenType struct {
	ID              TokenID
	Patterns        []string
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewMultiTokenType(id TokenID, patterns []string, ignoreCase bool, priority int, previousTokenId TokenID) *MultiTokenType {
	return &MultiTokenType{
		ID:              id,
		Patterns:        patterns,
		IgnoreCase:      ignoreCase,
		Priority:        priority,
		PreviousTokenID: previousTokenId,
	}
}

func (ptt *MultiTokenType) String() string {
	return ptt.ID.String()
}

func (ptt *MultiTokenType) GetID() TokenID {
	return ptt.ID
}

func (ptt *MultiTokenType) FindToken(s string, p Position) *Token {
	var found []string
	if ptt.IgnoreCase {
		s = strings.ToUpper(s)
	}
	for _, x := range ptt.Patterns {
		if strings.HasPrefix(s, x) {
			found = append(found, x)
		}
	}

	if len(found) == 0 {
		return nil
	}

	// Use longest pattern
	maxIdx := 0
	maxLen := len(found[0])
	for i, s := range found {
		if len(s) > maxLen {
			maxIdx = i
			maxLen = len(s)
		}
	}

	return &Token{
		Type:     ptt,
		Literal:  found[maxIdx],
		Position: p,
	}
}

func (ptt *MultiTokenType) GetPreviousTokenID() TokenID {
	return ptt.PreviousTokenID
}

func (ptt *MultiTokenType) GetPriority() int {
	return ptt.Priority
}

type RegexpTokenType struct {
	ID              TokenID
	Re              *regexp.Regexp
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewRegexpTokenType(id TokenID, re string, ignoreCase bool, priority int, previousTokenId TokenID) *RegexpTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &RegexpTokenType{
		ID:              id,
		Re:              regexp.MustCompile(re),
		IgnoreCase:      ignoreCase,
		Priority:        priority,
		PreviousTokenID: previousTokenId,
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
	ID              TokenID
	Re              *regexp.Regexp
	Parentheses     []rune
	IgnoreCase      bool
	Priority        int
	PreviousTokenID TokenID
}

func NewBalancedParenthesesTokenType(id TokenID, re string, open rune, close rune, ignoreCase bool, priority int, previousTokenId TokenID) *BalancedParenthesesTokenType {
	if !strings.HasPrefix(re, "^") {
		re = "^(?:" + re + ")"
	}
	return &BalancedParenthesesTokenType{
		ID:              id,
		Re:              regexp.MustCompile(re),
		Parentheses:     []rune{open, close},
		IgnoreCase:      ignoreCase,
		Priority:        priority,
		PreviousTokenID: previousTokenId,
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

func (btt *BalancedParenthesesTokenType) GetPreviousTokenID() TokenID {
	return btt.PreviousTokenID
}

func (btt *BalancedParenthesesTokenType) GetPriority() int {
	return btt.Priority
}

type Token struct {
	Type       TokenType
	Literal    string   // The string of matched.
	Submatches []string // Submatches of regular expression.
	Position   Position // Position of token.
}

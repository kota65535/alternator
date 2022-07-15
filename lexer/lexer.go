package lexer

import (
	"github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
)

const MaxInt = int(^uint(0) >> 1)

type Lexer struct {
	reader            io.Reader
	buf               string
	loadedLine        string
	nextPos           Position
	TokenTypes        []TokenType
	SkippedTokenTypes []TokenType
	prevTokens        []*Token
}

func NewLexer(reader io.Reader, tokenTypes []TokenType, skippedTokenTypes []TokenType) *Lexer {
	return &Lexer{
		reader:            reader,
		TokenTypes:        tokenTypes,
		SkippedTokenTypes: skippedTokenTypes,
	}
}

func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()
	l.consumeBuffer(t)
	return t, e
}

func (l *Lexer) Peek() (*Token, error) {

	// Ensure sorted by priority
	sort.SliceStable(l.TokenTypes, func(i, j int) bool {
		return l.TokenTypes[i].GetPriority() < l.TokenTypes[j].GetPriority()
	})

	// Find all matching tokens
	curPriority := MaxInt
	var found []*Token
	for _, tokenType := range l.TokenTypes {
		// Check priority.
		// If once found, all tokens with lower priority (higher value) are skipped
		if curPriority < tokenType.GetPriority() {
			continue
		}
		l.skipTokens()
		l.readBufIfNeed()
		if t := tokenType.FindToken(l.buf, l.nextPos); t != nil {
			found = append(found, t)
			curPriority = tokenType.GetPriority()
		}
	}

	if len(found) == 0 {
		if len(l.buf) > 0 {
			return nil, l.makeError()
		} else {
			return nil, nil
		}
	}

	for _, tokenType := range l.TokenTypes {
		l.skipTokens()
		l.readBufIfNeed()
		if t := tokenType.FindToken(l.buf, l.nextPos); t != nil {
			found = append(found, t)
		}
	}

	if len(found) == 0 {
		if len(l.buf) > 0 {
			return nil, l.makeError()
		} else {
			return nil, nil
		}
	}

	// Select the longest of the matched tokens
	maxIdx := 0
	maxLen := len(found[0].Literal)
	for i, t := range found {
		curLen := len(t.Literal)
		if maxLen < curLen {
			maxIdx = i
			maxLen = curLen
		}
	}

	//fmt.Fprintf(os.Stderr, "Token = %s, %d\n", found[maxIdx].Literal, found[maxIdx].Type.GetID())

	l.prevTokens = append(l.prevTokens, found[maxIdx])
	return found[maxIdx], nil
}

func (l *Lexer) readBufIfNeed() {
	if len(l.buf) < 1024 {
		buf := make([]byte, 2048)
		l.reader.Read(buf)
		l.buf += strings.TrimRight(string(buf), "\x00")
	}
}

func (l *Lexer) consumeBuffer(t *Token) {
	if t == nil {
		return
	}

	l.buf = l.buf[len(t.Literal):]

	l.nextPos = shiftPos(l.nextPos, t.Literal)

	if idx := strings.LastIndex(t.Literal, "\n"); idx >= 0 {
		l.loadedLine = t.Literal[idx+1:]
	} else {
		l.loadedLine += t.Literal
	}
}

func (l *Lexer) skipTokens() {
	for true {
		l.readBufIfNeed()

		skipped := false
		for _, t := range l.SkippedTokenTypes {
			if tok := t.FindToken(l.buf, l.nextPos); tok != nil {
				logrus.Debugf("skip %s", tok.Literal)
				l.consumeBuffer(tok)
				skipped = true
			}
		}
		if !skipped {
			break
		}
	}
}

func (l *Lexer) makeError() error {
	return UnknownTokenError{
		Literal:  l.buf,
		Position: l.nextPos,
	}
}

func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()

	if idx := strings.Index(l.buf, "\n"); idx >= 0 {
		return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
	} else {
		return l.loadedLine + l.buf
	}
}

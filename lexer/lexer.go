package lexer

import (
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
	LastToken         TokenType
}

func NewLexer(reader io.Reader, tokenTypes []TokenType, skippedTokenTypes []TokenType) *Lexer {
	return &Lexer{
		reader:            reader,
		TokenTypes:        tokenTypes,
		SkippedTokenTypes: skippedTokenTypes,
	}
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
	//for shift, _ := range l.buf {
	//	if l.SkippedTokenTypes != nil && l.SkippedTokenTypes.FindToken(l.buf[shift:], l.nextPos) != nil {
	//		return UnknownTokenError{
	//			Literal:  l.buf[:shift],
	//			Position: l.nextPos,
	//		}
	//	}
	//
	//	for _, tokenType := range l.TokenTypes {
	//		if tokenType.FindToken(l.buf[shift:], l.nextPos) != nil {
	//			return UnknownTokenError{
	//				Literal:  l.buf[:shift],
	//				Position: l.nextPos,
	//			}
	//		}
	//	}
	//}

	return UnknownTokenError{
		Literal:  l.buf,
		Position: l.nextPos,
	}
}

func (l *Lexer) Peek() (*Token, error) {
	// Ensure sorted by priority
	sort.SliceStable(l.TokenTypes, func(i, j int) bool {
		return l.TokenTypes[i].GetPriority() < l.TokenTypes[j].GetPriority()
	})

	curPriority := MaxInt
	var found []*Token
	for _, tokenType := range l.TokenTypes {
		// Check priority.
		// If once found, all tokens with lower priority (higher value) are skipped
		if curPriority < tokenType.GetPriority() {
			continue
		}
		// Check previous token constraints
		if prevTokenId := tokenType.GetPreviousTokenID(); prevTokenId > 0 && l.LastToken != nil && prevTokenId != l.LastToken.GetID() {
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

	l.LastToken = found[maxIdx].Type
	return found[maxIdx], nil
}

/*
Scan will get the first token in the buffer and remove it from the buffer.

This function using Lexer.Peek. Please read document of Peek.
*/
func (l *Lexer) Scan() (*Token, error) {
	t, e := l.Peek()

	l.consumeBuffer(t)

	return t, e
}

func (l *Lexer) GetLastLine() string {
	l.readBufIfNeed()

	if idx := strings.Index(l.buf, "\n"); idx >= 0 {
		return l.loadedLine + l.buf[:strings.Index(l.buf, "\n")]
	} else {
		return l.loadedLine + l.buf
	}
}

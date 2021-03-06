package lexer

import "fmt"

type UnknownTokenError struct {
	Literal  string
	Position Position
}

// Get error message as string.
func (se UnknownTokenError) Error() string {
	return fmt.Sprintf("%d:%d:UnknownTokenError: %#v", se.Position.Line+1, se.Position.Column+1, se.Literal)
}

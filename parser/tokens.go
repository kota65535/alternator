package parser

import "github.com/kota65535/alternator/lexer"

var Signs = map[int]string{
	lp:        "(",
	rp:        ")",
	lcb:       "{",
	rcb:       "}",
	comma:     ",",
	semicolon: ";",
	eq:        "=",
	dot:       ".",
	gt:        ">",
	gte:       ">=",
	lt:        "<",
	lte:       "<=",
	ne:        "!=",
	ne2:       "<>",
	nseq:      "<=>",
	tilde:     "~",
	and:       "&",
	and2:      "&&",
	or:        "|",
	or2:       "||",
	rshift:    "<<",
	lshift:    ">>",
	plus:      "+",
	minus:     "-",
	mult:      "*",
	div:       "/",
	mod:       "%",
	hat:       "^",
	excl:      "!",
	qstn:      "?",
}

var Literals = []lexer.TokenType{
	lexer.NewRegexpTokenType(BIT_STR, `[bB]'[01]+'`, 2),
	lexer.NewRegexpTokenType(BIT_NUM, `0[bB][01]+`, 2),
	lexer.NewRegexpTokenType(INT_NUM, `[0-9]+`, 2),
	lexer.NewRegexpTokenType(HEX_STR, `[xX]'[0-9A-F]+'`, 2),
	lexer.NewRegexpTokenType(HEX_NUM, `0x[0-9a-fA-F]+`, 2),
	lexer.NewRegexpTokenType(FLOAT_NUM, `[+-]?[0-9]+(\.[0-9]+)?(E[+-]?[0-9]+)?`, 2),
	lexer.NewRegexpTokenType(STRING, `'(?:[^'\\]|.)*?'`, 2),
	lexer.NewRegexpTokenType(IDENTIFIER, `[a-zA-Z_][a-zA-Z0-9_]*`, 2),
	lexer.NewRegexpTokenType(LOCAL_VAR, "@[a-zA-Z_][a-zA-Z0-9_]*", 2),
	lexer.NewRegexpTokenType(LOCAL_VAR, "@`[a-zA-Z_][a-zA-Z0-9_]*`", 2),
	lexer.NewRegexpTokenType(GLOBAL_VAR, "@@(GLOBAL\\.|SESSION\\.)?[a-zA-Z_][a-zA-Z0-9_]*", 2),
	lexer.NewRegexpTokenType(GLOBAL_VAR, "@@`(GLOBAL\\.|SESSION\\.)?[a-zA-Z_][a-zA-Z0-9_]*`", 2),
	lexer.NewRegexpTokenType(QUOTED_IDENTIFIER, "`([a-zA-Z_][a-zA-Z0-9_]*)`", 2),
}

var Skipped = []string{
	`\s`,                // space
	`#.*\n`,             // single line comment
	`--.*\n`,            // single line comment
	`(?s)/\*[^!].*?\*/`, // multiple line comment
	`/\*!\d{5}`,         // begin of comment with MySQL extensions
	`\*/`,               // end of comment with MySQL extensions
}

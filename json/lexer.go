package json

import (
	"fmt"
)

type tokenType string

const (
	lBracketType = tokenType("[")
	rBracketType = tokenType("]")
	lBraceType   = tokenType("{")
	rBraceType   = tokenType("}")
	colonType    = tokenType(":")
	commaType    = tokenType(",")
	stringType   = tokenType("string")
	numberType   = tokenType("number")
	booleanType  = tokenType("bool")
	nullType     = tokenType("null")
	eofType      = tokenType("EOF")
)

type (
	token struct {
		tokenType tokenType
		char      string
	}

	lexer struct {
		str []rune
		idx int
	}

	ErrInvalidToken struct {
		expected string
		got      string
	}
)

func (e ErrInvalidToken) Error() string {
	return fmt.Sprintf("invalid token is found: expected %s but got %s", e.expected, e.got)
}

func newLexer(str string) *lexer {
	return &lexer{str: []rune(str), idx: 0}
}

func (l *lexer) scan() (*token, error) {
	if len(l.str) == l.idx {
		return &token{tokenType: eofType, char: ""}, nil
	}

	char := l.next()
	for {
		switch char {
		case "\n", "\t", " ":
			char = l.next()
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			return l.numberType(char)
		case "t":
			return l.booleanType("true")
		case "f":
			return l.booleanType("false")
		case "n":
			return l.nullType()
		case `"`:
			return l.stringType()
		case "[":
			return &token{tokenType: lBracketType, char: char}, nil
		case "]":
			return &token{tokenType: rBracketType, char: char}, nil
		case "{":
			return &token{tokenType: lBraceType, char: char}, nil
		case "}":
			return &token{tokenType: rBraceType, char: char}, nil
		case ":":
			return &token{tokenType: colonType, char: char}, nil
		case ",":
			return &token{tokenType: commaType, char: char}, nil
		default:
			return nil, ErrInvalidToken{expected: "number, t, f, n, \", [, ], {, }, : or ,", got: char}
		}
	}
}

func (l *lexer) numberType(f string) (*token, error) {
	num := f
	isDecimal := false
	endWithDot := false
	for isNum := true; isNum; {
		next := l.peek()
		switch next {
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			endWithDot = false
			num = fmt.Sprintf("%s%s", num, l.next())
		case ".":
			if isDecimal {
				return nil, ErrInvalidToken{expected: "number", got: next}
			}

			isDecimal = true
			endWithDot = true
			num = fmt.Sprintf("%s%s", num, l.next())
		default:
			isNum = false
			if endWithDot {
				return nil, ErrInvalidToken{expected: "number or .", got: next}
			}
		}
	}
	return &token{tokenType: numberType, char: num}, nil
}

func (l *lexer) booleanType(b string) (*token, error) {
	for i, c := range b {
		if i == 0 {
			continue
		}

		if l.next() != string(c) {
			return nil, ErrInvalidToken{expected: string(c), got: l.next()}
		}
	}

	return &token{tokenType: booleanType, char: b}, nil
}

func (l *lexer) nullType() (*token, error) {
	for i, c := range "null" {
		if i == 0 {
			continue
		}

		if l.next() != string(c) {
			return nil, ErrInvalidToken{expected: string(c), got: l.next()}
		}
	}

	return &token{tokenType: nullType, char: "null"}, nil
}

func (l *lexer) stringType() (*token, error) {
	var char string
	for {
		n := l.next()
		switch n {
		case `\`:
			nn := l.next()
			char = fmt.Sprintf("%s%s%s", char, n, nn)
		case `"`:
			return &token{tokenType: stringType, char: char}, nil
		case "":
			return nil, ErrInvalidToken{expected: "other than EOF", got: "EOF"}
		default:
			char = fmt.Sprintf("%s%s", char, n)
		}
	}
}

func (l *lexer) next() string {
	if len([]rune(l.str)) == l.idx {
		return ""
	}

	char := []rune(l.str)[l.idx]
	l.idx++
	return string(char)
}

func (l *lexer) peek() string {
	if len([]rune(l.str)) == l.idx {
		return ""
	}

	return string([]rune(l.str)[l.idx])
}

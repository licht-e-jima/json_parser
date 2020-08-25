package json

import (
	"errors"
	"fmt"
	"strconv"
)

type (
	parser struct {
		lexer        *lexer
		parsingToken *token
	}

	ErrUnexpectedToken struct {
		expected string
		got      string
	}
)

var (
	ErrInvalidNode error = errors.New("unexpected node is found")
)

func (e ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token is found: expected %s but got %s", e.expected, e.got)
}

func newParser(str string) *parser {
	lexer := newLexer(str)
	return &parser{lexer: lexer, parsingToken: nil}
}

// ll1 parse string into a tree using LL1 method.
//
// json -> entity eofType
// entity -> numberType | stringType | booleanType | nullType | array | structure
// array -> lBracketType entity entityCircuit rBracketType | lBracketType entity rBracketType | lBracketType rBracketType
// entityCircuit -> commaType entity entityCircuit | commaType entity
// structure -> lBraceType keyValue keyValueCircuit rBraceType | lBraceType keyValue rBracketType | lBraceType rBracketType
// keyValue -> stringType colonType entity
// keyValueCircuit -> commaType keyValue keyValueCircuit | commaType keyValue
func (p *parser) ll1() (node, error) {
	if err := p.move(); err != nil {
		return nil, err
	}
	tree, err := p.json()
	if err != nil {
		return nil, err
	}
	return tree, nil
}

// json -> entity eofType
func (p *parser) json() (node, error) {
	n, err := p.entity()
	if err != nil {
		return nil, err
	}
	if err := p.match(eofType); err != nil {
		return nil, err
	}
	return n, nil
}

// entity -> numberType | stringType | booleanType | nullType | array | structure
func (p *parser) entity() (node, error) {
	char := p.parsingToken.char
	tt := p.parsingToken.tokenType
	if tt == numberType {
		if err := p.match(numberType); err != nil {
			return nil, err
		}

		return &numberNode{char: char}, nil
	} else if tt == stringType {
		if err := p.match(stringType); err != nil {
			return nil, err
		}

		return &stringNode{char: char}, nil
	} else if tt == booleanType {
		b, err := strconv.ParseBool(char)
		if err != nil {
			return nil, err
		}

		if err := p.match(booleanType); err != nil {
			return nil, err
		}

		return &booleanNode{boolean: b}, nil
	} else if tt == nullType {
		if err := p.match(nullType); err != nil {
			return nil, err
		}

		return &nullNode{}, nil
	} else if p.isArray() {
		return p.array()
	} else if p.isStructure() {
		return p.structure()
	} else {
		return nil, ErrUnexpectedToken{expected: "number, string, boolean, null, [ or {", got: char}
	}
}

// array -> lBracketType entity entityCircuit rBracketType | lBracketType entity rBracketType | lBracketType
// entityCircuit -> commaType entity entityCircuit | commaType entity
func (p *parser) array() (node, error) {
	if err := p.match(lBracketType); err != nil {
		return nil, err
	}

	if p.parsingToken.tokenType == rBracketType {
		if err := p.match(rBracketType); err != nil {
			return nil, err
		}

		return &arrayNode{}, nil
	}

	elements := []node{}
	e, err := p.entity()
	if err != nil {
		return nil, err
	}
	elements = append(elements, e)
	for t := p.parsingToken; t.tokenType == commaType; t = p.parsingToken {
		if err := p.match(commaType); err != nil {
			return nil, err
		}

		e, err = p.entity()
		if err != nil {
			return nil, err
		}
		elements = append(elements, e)
	}

	if err := p.match(rBracketType); err != nil {
		return nil, err
	}
	return &arrayNode{elements: elements}, nil
}

func (p *parser) isArray() bool {
	return p.parsingToken.tokenType == lBracketType
}

// structure -> lBraceType keyValue keyValueCircuit rBraceType | lBraceType keyValue rBracketType | lBraceType rBracketType
// keyValueCircuit -> commaType keyValue keyValueCircuit | commaType keyValue
func (p *parser) structure() (node, error) {
	if err := p.match(lBraceType); err != nil {
		return nil, err
	}

	if p.parsingToken.tokenType == rBraceType {
		if err := p.match(rBraceType); err != nil {
			return nil, err
		}
		return &structureNode{}, nil
	}

	values := map[string]node{}
	k, v, err := p.keyValue()
	if err != nil {
		return nil, err
	}
	values[k] = v
	for t := p.parsingToken; t.tokenType == commaType; t = p.parsingToken {
		if err := p.match(commaType); err != nil {
			return nil, err
		}

		k, v, err = p.keyValue()
		if err != nil {
			return nil, err
		}
		values[k] = v
	}

	if err := p.match(rBraceType); err != nil {
		return nil, err
	}

	return &structureNode{values: values}, nil
}

func (p *parser) isStructure() bool {
	return p.parsingToken.tokenType == lBraceType
}

// keyValue -> stringType colonType entity
func (p *parser) keyValue() (string, node, error) {
	key := p.parsingToken.char
	if err := p.match(stringType); err != nil {
		return "", nil, err
	}

	if err := p.match(colonType); err != nil {
		return "", nil, err
	}

	val, err := p.entity()
	if err != nil {
		return "", nil, err
	}

	return key, val, nil
}

func (p *parser) match(tokenType tokenType) error {
	if p.parsingToken.tokenType != tokenType {
		return ErrUnexpectedToken{expected: string(tokenType), got: p.parsingToken.char}
	}

	return p.move()
}

func (p *parser) move() error {
	token, err := p.lexer.scan()
	if err != nil {
		return err
	}
	p.parsingToken = token
	return nil
}

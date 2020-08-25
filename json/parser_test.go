package json

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var jsonStr = `{
	"string": "foo",
	"rune": "ほげ",
	"number": 1,
	"decimal": 1.1,
	"true": true,
	"false": false,
	"null": null,
	"array": [
		"bar",
		"ふが",
		2,
		2.2,
		true,
		false,
		null,
		[],
		{}
	]
}`

func TestParser_ll1(t *testing.T) {
	type expected struct {
		node node
		err  error
	}
	tests := []struct {
		name     string
		str      string
		expected expected
	}{
		{
			name: "正常",
			str:  jsonStr,
			expected: expected{
				node: &structureNode{
					values: map[string]node{
						"string":  &stringNode{char: "foo"},
						"rune":    &stringNode{char: "ほげ"},
						"number":  &numberNode{char: "1"},
						"decimal": &numberNode{char: "1.1"},
						"true":    &booleanNode{boolean: true},
						"false":   &booleanNode{boolean: false},
						"null":    &nullNode{},
						"array": &arrayNode{
							elements: []node{
								&stringNode{char: "bar"},
								&stringNode{char: "ふが"},
								&numberNode{char: "2"},
								&numberNode{char: "2.2"},
								&booleanNode{boolean: true},
								&booleanNode{boolean: false},
								&nullNode{},
								&arrayNode{},
								&structureNode{},
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "無効なJSON",
			str:  fmt.Sprintf("%s,", jsonStr),
			expected: expected{
				node: nil,
				err:  ErrUnexpectedToken{expected: "EOF", got: ","},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			p := newParser(test.str)
			tree, err := p.ll1()
			assert.Equal(t, test.expected.node, tree)
			assert.Equal(t, test.expected.err, err)
		})
	}

}

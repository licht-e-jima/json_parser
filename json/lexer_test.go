package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer_scan(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		expectChar string
		tokenType  tokenType
		err        error
	}{
		{
			name:      "整数",
			str:       "0123456789",
			tokenType: numberType,
			err:       nil,
		},
		{
			name:      "小数",
			str:       "1.1",
			tokenType: numberType,
			err:       nil,
		},
		{
			name:       "文字列",
			str:        `"abcde\"fgh"`,
			expectChar: `abcde\"fgh`,
			tokenType:  stringType,
			err:        nil,
		},
		{
			name:       "文字列(multi bytes)",
			str:        `"あいうえお"`,
			expectChar: `あいうえお`,
			tokenType:  stringType,
			err:        nil,
		},
		{
			name:      "true",
			str:       "true",
			tokenType: booleanType,
			err:       nil,
		},
		{
			name:      "false",
			str:       "false",
			tokenType: booleanType,
			err:       nil,
		},
		{
			name:      "null",
			str:       "null",
			tokenType: nullType,
			err:       nil,
		},
		{
			name:      "{",
			str:       "{",
			tokenType: lBraceType,
			err:       nil,
		},
		{
			name:      "}",
			str:       "}",
			tokenType: rBraceType,
			err:       nil,
		},
		{
			name:      "[",
			str:       "[",
			tokenType: lBracketType,
			err:       nil,
		},
		{
			name:      "]",
			str:       "]",
			tokenType: rBracketType,
			err:       nil,
		},
		{
			name:      ":",
			str:       ":",
			tokenType: colonType,
			err:       nil,
		},
		{
			name: "無効なjson",
			str:  "abc",
			err:  ErrInvalidToken,
		},
		{
			name: "無効な数字",
			str:  "1.",
			err:  ErrInvalidToken,
		},
		{
			name: "無効な数字",
			str:  ".1",
			err:  ErrInvalidToken,
		},
		{
			name: "無効な数字",
			str:  "1..1",
			err:  ErrInvalidToken,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			l := newLexer(test.str)
			if test.err == nil {
				var (
					token *token
					err   error
				)
				for token, err = l.scan(); err == nil && token.tokenType != eofType; token, err = l.scan() {
					assert.Equal(t, test.tokenType, token.tokenType)
					var expect string
					if test.expectChar != "" {
						expect = test.expectChar
					} else {
						expect = test.str
					}
					assert.Equal(t, expect, token.char)
					assert.NoError(t, err)
				}
				assert.NoError(t, err)
			} else {
				_, err := l.scan()
				assert.Equal(t, test.err, err)
			}
		})
	}
}

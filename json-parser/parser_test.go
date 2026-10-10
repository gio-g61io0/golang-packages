package jsonparser_test

import (
	jsonparser "personal-http-server/json-parser"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestNumericLexer(t *testing.T){
	input := "12345"
	lexer := jsonparser.NewLexer(input)
	token, _, err := lexer.ConsumeAsNumeric(0,input)

	assert.Equal(t, jsonparser.Token{
		Line: 0,
		Value: "12345",
	}, token)
	assert.NoError(t, err)
}

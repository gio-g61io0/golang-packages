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
func TestStringLexer(t *testing.T){
	input := "\"\""
	lexer := jsonparser.NewLexer(input)
	token, _, err := lexer.ConsumeAsString(input)

	assert.Equal(t, jsonparser.Token{
		Line: 2,
		Value: "\"\"",
		Type: jsonparser.StringToken,
	}, token)
	assert.NoError(t, err)
	input = "\"Gio Gonzales\""
	token, _, err = lexer.ConsumeAsString(input)
	assert.Equal(t, jsonparser.Token{
		Line: 14,
		Value: "\"Gio Gonzales\"",
		Type: jsonparser.StringToken,
	}, token)
	assert.NoError(t, err)
}

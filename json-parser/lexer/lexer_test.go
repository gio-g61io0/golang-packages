package lexer_test

import (
	lexer "personal-http-server/json-parser/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestNumericLexer(t *testing.T){
	input := "12345"
	lex:= lexer.NewLexer(input)
	token, _, err := lex.ConsumeAsNumeric(input)

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: "12345",
	}, token)
	assert.NoError(t, err)
}
func TestStringLexer(t *testing.T){
	input := "\"\""
	lex:= lexer.NewLexer(input)
	token, parsed, err := lex.ConsumeAsString(input)

	assert.Equal(t, lexer.Token{
		Line: 2,
		Value: "\"\"",
		Type: lexer.StringToken,
	}, token)

	assert.NoError(t, err)
	input = "\"Gio Gonzales\","
	lex = lexer.NewLexer(input)

	token, parsed, err = lex.ConsumeAsString(input)
	lex.UpdateInput(parsed)
	assert.Equal(t, lexer.Token{
		Line: 14,
		Value: "\"Gio Gonzales\"",
		Type: lexer.StringToken,
	}, token)
	assert.Equal(t, ",", lex.Input)

	assert.NoError(t, err)

	input = "\"name\":\"Cj\","
	lex = lexer.NewLexer(input)

	token, parsed, err = lex.ConsumeAsString(input)
	lex.UpdateInput(parsed)
	assert.Equal(t, lexer.Token{
		Line: 6,
		Value: "\"name\"",
		Type: lexer.StringToken,
	}, token)

	assert.Equal(t, ":\"Cj\",", lex.Input)
	assert.NoError(t, err)
}
func TestEmptyStringLexer(t *testing.T){

	input := ""
	lex := lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
	}, token)
	assert.Equal(t, 0, parsed)
	assert.Equal(t, "", remaining)
	assert.Error(t, err)
}

func TestCommaLexer(t *testing.T){
	input := ","
	lex:= lexer.NewLexer(input)
	token, remaining, _, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: ",",
		Type: lexer.CommaToken,
	}, token)
	assert.Equal(t, "", remaining)
	assert.NoError(t, err)

}
func TestColonLexer(t *testing.T){
	input := ":\"Gio\""
	lex:= lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: ":",
		Type: lexer.ColonToken,
	}, token)
	assert.Equal(t, "\"Gio\"", remaining)
	assert.Equal(t, 1, parsed)
	assert.NoError(t, err)
}
func TestObjectOpenerLexer(t *testing.T){
	input := "{[]}"
	lex:= lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: "{",
		Type: lexer.ObjectOpenerToken,
	}, token)
	assert.Equal(t, "[]}", remaining)
	assert.Equal(t, 1, parsed)
	assert.NoError(t, err)
}
func TestObjectCloserLexer(t *testing.T){
	input := "}}}"
	lex:= lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: "}",
		Type: lexer.ObjectCloserToken,
	}, token)
	assert.Equal(t, "}}", remaining)
	assert.Equal(t, 1, parsed)
	assert.NoError(t, err)
}
func TestArrayOpenerLexer(t *testing.T){
	input := "[[]]"
	lex:= lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: "[",
		Type: lexer.ArrayOpenerToken,
	}, token)
	assert.Equal(t, "[]]", remaining)
	assert.Equal(t, 1, parsed)
	assert.NoError(t, err)
}
func TestArrayCloserLexer(t *testing.T){
	input := "]][]"
	lex:= lexer.NewLexer(input)
	token, remaining, parsed, err := lex.Step()

	assert.Equal(t, lexer.Token{
		Line: 0,
		Value: "]",
		Type: lexer.ArrayCloserToken,
	}, token)
	assert.Equal(t, "][]", remaining)
	assert.Equal(t, 1, parsed)
	assert.NoError(t, err)
}

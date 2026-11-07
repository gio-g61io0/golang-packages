package jsonparser

import (
	"fmt"
	"log/slog"
	"personal-http-server/json-parser/lexer"
	"strings"
)


type Parser struct{
	input string
	LineNumber int
	lexer *lexer.Lexer
}

func NewParser(input string) *Parser{
	return &Parser{
		input: input,
		LineNumber: 0,
		lexer: lexer.NewLexer(input),
	}
}

func (parser *Parser)ParseStringToken(stringToken lexer.Token) error{
	if stringToken.Type != lexer.StringToken{
		return fmt.Errorf("Cannot parse non string token")
	}
	return nil
}
//1. Let us try to parse a key:value first
func (parser *Parser)MatchToken(token lexer.Token) error{
	switch token.Type{
	case lexer.StringToken:
	parser.ParseStringToken(token)
	default:
	return nil
	}
	return nil

}

func (parser *Parser)Parse() error{
	var token lexer.Token
	var parsed int
	var err error

	//iterates every line and for every line, we iteratively produce a token 
	
	return nil
}

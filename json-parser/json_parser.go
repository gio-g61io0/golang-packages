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
	for line := range strings.Lines(parser.input){

		jsonString := line

		for  len(jsonString) != 0 {
			token, jsonString, parsed, err = parser.lexer.Step(jsonString)

			slog.Info("Remaning string input", "remaining", jsonString)
			if err != nil{
				slog.Error("Something went wrong producing a token", "error", err)
				return err
			}
			if parsed == 0{
				slog.Info("No token is produced for the remaining sequence of chars in a line")
				break
			}
		}
		parser.LineNumber  += 1
	}
	return nil
}

package jsonparser

import (
	"fmt"
	"strings"
	"unicode"
)
const NEWLINE = "\n"


type Token struct{
	Line int
	Value string
}
type Lexer struct{
	Input string //Sequence of input string 
	Tokens []Token
	Position int
}

//Will try to figure out if we can produce token in this line
//This assumes that prior to calling this function, the line was checked that it contains atleast numeric character
func (lexer *Lexer)ConsumeAsNumeric(starting int , line string) (Token, int,  error) {
	//Gives us the idx of the character within this line that is not numeric
	idx := strings.IndexFunc(line, func(char rune) bool{
		return !unicode.IsDigit(char)
	})

	if idx == -1{
		return Token{Value: line, Line: 0}, idx, nil
	}
	return Token{
		Value: line[:idx],
		Line: 0,
	}, idx, nil

}

func (lexer *Lexer)Lex(line string) (Token, string, error) {
	//Clean the string first
	cleanedLine := strings.TrimSpace(line)

	//I think I will only need to check the first character?
	for idx, char := range cleanedLine {
		if unicode.IsDigit(char){
			numericToken, parsedChars, err := lexer.ConsumeAsNumeric(idx, cleanedLine[idx:])
			if err != nil {

			}
			return numericToken, line[parsedChars:], nil
			
		}
		switch char {
		case '"':
			return Token{}, "", nil
		}
	}
	return Token{}, "", nil
}

func (lexer *Lexer)Analyze() ([]Token, string, error) {
	for{
		relativePosition := strings.Index(lexer.Input[lexer.Position:], NEWLINE)

		if relativePosition  == -1{
			break;
		}

		position := lexer.Position + relativePosition
		fmt.Printf("%s\n", lexer.Input[lexer.Position:position])
		lexer.Position = position + len(NEWLINE)
	}
	return lexer.Tokens, "", nil
}
func NewLexer(inputString string) *Lexer{
	return &Lexer{
		Position: 0,
		Input: inputString,
		Tokens: []Token{},
	}
}



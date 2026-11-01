package jsonparser

import (
	"fmt"
	"strings"
	"unicode"
)
const NEWLINE = "\n"

type TokenType int

const (
	NumericToken TokenType= iota
	CommaToken 
	StringToken
	ColonToken
	ObjectOpenerToken
	ObjectCloserToken
	ArrayOpenerToken
	ArrayCloserToken
)


type Token struct{
	Line int
	Value string
	Type TokenType
}
type Lexer struct{
	Input string //Sequence of input string 
	Tokens []Token
	Position int
	Line int
	CharIdx int
}

//Will try to figure out if we can produce token in this line
//This assumes that prior to calling this function, the line was checked that it contains atleast numeric character
func (lexer *Lexer)ConsumeAsNumeric(starting int , line string) (Token, int,  error) {
	//Gives us the idx of the character within this line that is not numeric
	idx := strings.IndexFunc(line, func(char rune) bool{
		return !unicode.IsDigit(char)
	})

	if idx == -1{
		return Token{Value: line, Line: 0, Type: NumericToken}, idx + 1, nil
	}
	return Token{
		Value: line[:idx],
		Line: 0,
	}, idx + 1, nil

}
// func (lexer *Lexer)Lex(line string) (Token, string, error) {
//
// }


func (lexer *Lexer)updateCharIdx(idx int) {
	lexer.CharIdx += idx
}

func (lexer *Lexer)ConsumeAsString(data string) (Token, int, error) {
	idx := strings.IndexFunc(data[1:] , func(char rune) bool{
		return char == '"'
	})

	idx += 1

	if idx == -1{
		return Token{}, 0, fmt.Errorf("Cannot be tokenized as string")
	}
	return Token{Value: data[:idx + 1], Line: idx + 1, Type:StringToken }, idx + 1, nil
}

//This will return a single token that it can see
func (lexer *Lexer)Step(line string) (Token, string, int,  error) {
	//Clean the string first
	cleanedLine := strings.TrimSpace(line)
	
	if len(line) == 0{
		return Token{}, "", 0, fmt.Errorf("Cannot parse empty string")
	}

	//I think I will only need to check the first character?
	runes := []rune(cleanedLine)
	char := runes[0]

	if unicode.IsDigit(char){
		numericToken, parsedChars, err := lexer.ConsumeAsNumeric(0, cleanedLine)
		if err != nil {
			return Token{}, cleanedLine, 0, err
		}
		lexer.updateCharIdx(parsedChars)
		return numericToken, cleanedLine[parsedChars:], parsedChars, nil
	}
	switch char {
	case '"':
		stringToken, parsedChars, err := lexer.ConsumeAsString(cleanedLine)
		if err != nil {
			return Token{}, line, 0, err
		}
		return stringToken, cleanedLine[parsedChars:], parsedChars, nil
	case ':':
		return Token{
			Line: 0,
			Value: ":",
			Type: ColonToken,
		}, line[1:], 1, nil
	case ',':
		return Token{
			Line: 0,
			Value: ",",
			Type: CommaToken,
		}, line[1:], 1, nil
	case '{':
		return Token{
			Line: 0,
			Value: "{",
			Type: ObjectOpenerToken,
		}, line[1:], 1, nil
	case '}':
		return Token{
			Line: 0,
			Value: "}",
			Type: ObjectCloserToken,
		}, line[1:], 1, nil
	case '[':
		return Token{
			Line: 0,
			Value: "[",
			Type: ArrayOpenerToken,
		}, line[1:], 1, nil
	case ']':
		return Token{
			Line: 0,
			Value: "]",
			Type: ArrayCloserToken,
		}, line[1:], 1, nil
	default:
		return Token{}, line, 0, nil
	}
}

func (lexer *Lexer)Analyze() ([]Token, string, error) {
	for{
		relativePosition := strings.Index(lexer.Input[lexer.Position:], NEWLINE)

		if relativePosition  == -1{
			break;
		}

		lexer.Line += 1 //A new Line
		position := lexer.Position + relativePosition
		fmt.Printf("%s\n", lexer.Input[lexer.Position:position])
		lexer.Position = position + len(NEWLINE)
	}
	return lexer.Tokens, "", nil
}
func NewLexer(inputString string) *Lexer{
	return &Lexer{
		Position: 0,
		CharIdx:0,
		Line: 0,
		Input: inputString,
		Tokens: []Token{},
	}
}



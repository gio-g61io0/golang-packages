package lexer

import (
	"fmt"
	"log/slog"
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
	EOFToken
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
}

//Will try to figure out if we can produce token in this line
//This assumes that prior to calling this function, the line was checked that it contains atleast numeric character
func (lexer *Lexer)ConsumeAsNumeric(line string) (Token, int,  error) {
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

func (lexer *Lexer)ConsumeAsString(data string) (Token, int, error) {
	idx := strings.IndexFunc(data[1:] , func(char rune) bool{
		return char == '"'
	})
	if idx == -1{
		return Token{}, 0, fmt.Errorf("Cannot be tokenized as string")
	}
	idx += 1
	return Token{Value: data[:idx + 1], Line: idx + 1, Type:StringToken }, idx + 1, nil
}

func (lexer *Lexer)UpdatePosition(position int) {
	lexer.Position = position
}
func (lexer *Lexer)UpdateLine(line int) {
	lexer.Line = line
}
func (lexer *Lexer)UpdateInput(parsedChars int) {
	lexer.Input = lexer.Input[parsedChars:]
}

func (lexer *Lexer)GoNewLine(){
	lexer.Input = strings.TrimPrefix(lexer.Input, "\n\t")
}
func (lexer *Lexer)CleanInputPrefix(){
	lexer.Input = strings.TrimLeftFunc(lexer.Input, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}
//This will return a single token that it can see
func (lexer *Lexer)Step() (Token, string, int,  error) {
	//Clean the string first
	var line string

	slog.Info("Remaning Input", "#Input", lexer.Input)
	if len(lexer.Input) == 0{
		return Token{}, "", 0, fmt.Errorf("Cannot produce token. No remaining input")
	}

	//This removes whitespaces in prio
	lexer.CleanInputPrefix()

	newLineIdx := strings.Index(lexer.Input, NEWLINE)
	
	if newLineIdx  == -1{
		//Might be the end line of the sequence of characters
		slog.Warn("No newline for the remaining input string")
		line = lexer.Input
	//Might be that we should go to the next line
	}else{
		line = lexer.Input[:newLineIdx + 1] //The whole line including the newline
	}

	slog.Info("Uncleaned Line", "Line", line)
	//I think I will only need to check the first character?
	cleanedLine := strings.TrimSpace(line)
	slog.Info("Cleaned Line", "Line", cleanedLine)


	//Nothing is left to lex
	if len(cleanedLine) == 0{
		return Token{
			Line:0,
			Value: "",
			Type: EOFToken,
		}, "", 0, nil
	}

	runes := []rune(cleanedLine)
	char := runes[0]


	if unicode.IsDigit(char){
		numericToken, parsedChars, err := lexer.ConsumeAsNumeric(cleanedLine)
		if err != nil {
			return Token{}, cleanedLine, 0, err
		}
		lexer.UpdateInput(parsedChars)
		return numericToken, cleanedLine[parsedChars:], parsedChars, nil
	}
	switch char {
	case '"':
		stringToken, parsedChars, err := lexer.ConsumeAsString(cleanedLine)
		if err != nil {
			return Token{}, line, 0, err
		}
		slog.Info("Number of characters parsed", "#Parsed", parsedChars)
		lexer.UpdateInput(parsedChars)
		return stringToken, cleanedLine[parsedChars:], parsedChars, nil
	case ':':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: ":",
			Type: ColonToken,
		}, line[1:], 1, nil
	case ',':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: ",",
			Type: CommaToken,
		}, line[1:], 1, nil
	case '{':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: "{",
			Type: ObjectOpenerToken,
		}, line[1:], 1, nil
	case '}':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: "}",
			Type: ObjectCloserToken,
		}, line[1:], 1, nil
	case '[':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: "[",
			Type: ArrayOpenerToken,
		}, line[1:], 1, nil
	case ']':
		lexer.UpdateInput(1)
		return Token{
			Line: 0,
			Value: "]",
			Type: ArrayCloserToken,
		}, line[1:], 1, nil
	default:
		//There should an error here
		slog.Error("Invalid character", "Char", char)
		return Token{}, line, 0, fmt.Errorf("Invalid character")
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
		Line: 0,
		Input: inputString,
		Tokens: []Token{},
	}
}



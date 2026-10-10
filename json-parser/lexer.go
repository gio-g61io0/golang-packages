package jsonparser


type Token struct{
	Line int
	Value string
}
type Lexer struct{
	Input string //Sequence of input string 
	Tokens []Token
}

func (lexer *Lexer)Analyze() ([]Token, string, error) {

	return lexer.Tokens, "", nil
}
func NewLexer(inputString string) *Lexer{
	return &Lexer{
		Input: inputString,
		Tokens: []Token{},
	}
}



package main

import (
	"log/slog"
	// jsonparser "personal-http-server/json-parser"
	lexer "personal-http-server/json-parser/lexer"
)


func main(){


	jsonString := `"name":"Gio",
	"Age": 12,
	"hobbies":["skateboarding", "surfing", "jogging", "drumming", "programming", "problem solving"],
	"address1234": "This is an address at San Vicente, Butuan City",
	"friends":
	[
		{
		"name":"Cj",
		"Age": 12,
		"hobbies":["skateboarding", "surfing", "jogging", "drumming", "programming", "problem solving"],
		"address1234": "This is an address at San Vicente, Butuan City",
		},
		{
		"name":"RC",
		"Age": 12,
		"hobbies":["skateboarding", "surfing", "jogging", "drumming", "programming", "problem solving"],
		"address1234": "This is an address at San Vicente, Butuan City",
		},
		{
		"name":"Sean",
		"Age": 12,
		"hobbies":["skateboarding", "surfing", "jogging", "drumming", "programming", "problem solving"],
		"address1234": "This is an address at San Vicente, Butuan City",
		},
		{
		"name":"Errol",
		"Age": 12,
		"hobbies":["skateboarding", "surfing", "jogging", "drumming", "programming", "problem solving"],
		"address1234": "This is an address at San Vicente, Butuan City",
		},
	],
	"lover":{
		"name":"Mary Soliva",
		"Age": 12,
		"hobbies":["drawing", "surfing", "jogging", "stitching", "teaching", "problem solving"],
		"address1234": "This is an address at San Vicente, Butuan City",
		}
	`
	var token lexer.Token
	// var parsed int
	var err error

	parsedToken := []lexer.Token{}
	lex := lexer.NewLexer(jsonString)

	for err == nil {
		token, jsonString, _, err = lex.Step()
		slog.Info("Token lexed", "Token", token)
		parsedToken = append(parsedToken, token)
	}

}

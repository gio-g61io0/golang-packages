package main

import (
	"fmt"
	"log/slog"
	jsonparser "personal-http-server/json-parser"
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
	var token jsonparser.Token
	var parsed int
	var err error

	parsedToken := []jsonparser.Token{}
	lexer := jsonparser.NewLexer(jsonString)
	for {

		token, jsonString, parsed, err = lexer.Step(jsonString)

		slog.Info("Remaning string input", "remaining", jsonString)
		
		
		if err != nil || parsed == 0 {
			slog.Error("Something went wrong", "error", err)
			break
		}
		parsedToken = append(parsedToken, token)
	}

	fmt.Println("Successfully Parsed Input string")
	fmt.Printf("Tokens analyzed %v", parsedToken)


}

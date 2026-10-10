package main

import (
	"fmt"
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
	lexer := jsonparser.NewLexer(jsonString)
	tokens, remaining, err := lexer.Analyze()

	if err != nil{
		panic(fmt.Sprintf("An error occured during lexical analysis %s\n Remaining strings %s ", err, remaining))
	}
	fmt.Println("Successfully Parsed Input string")
	fmt.Printf("Tokens analyzed %v", tokens)


}

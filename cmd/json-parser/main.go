package main

import (
	"log/slog"
	// jsonparser "personal-http-server/json-parser"
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
	if parser:= jsonparser.NewParser(jsonString); parser != nil{
		err := parser.Parse()
		if err != nil{
			slog.Error("Something went wrong parsing input json string")
		}
	}



}

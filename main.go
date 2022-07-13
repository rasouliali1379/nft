package main

import (
	"fmt"
	"html"
	"nft/cmd"

	_ "nft/docs"

	"github.com/logrusorgru/aurora"
)

// @title    Nft API Documentation
// @version  1.0

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Description for what is this security definition being used
func main() {
	if err := cmd.Runner.RootCmd().Execute(); err != nil {
		fmt.Printf("\n %v Failed to run command: %v %v\n\n ", aurora.White(html.UnescapeString("&#x274C;")), err, aurora.White(html.UnescapeString("&#x274C;")))
	}
}

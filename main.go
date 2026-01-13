package main

import (
	"fmt"
	"os"

	"rahu/analyser"
	"rahu/parser"
	"rahu/utils"
)

func main() {
	input := utils.ParseFile("test.py")
	fmt.Printf("Going to parse input string\n%s\n\n\n", input)
	p := parser.New(input) // Pass the same input
	module := p.Parse()
	utils.PrintAST(os.Stdout, module)

	fmt.Println("\nFirst pass, scopes: ")
	global := analyser.BuildScopes(module)
	utils.PrintScopes(os.Stdout, global)

	fmt.Println("\nSecond pass, resolution: ")

	errors, resolved := analyser.Resolve(module, global)

	if len(errors) == 0 {
		fmt.Println("No semantic errors")
	} else {
		for _, err := range errors {
			fmt.Printf("Error at %v: %s\n", err.Span, err.Msg)
		}
	}

	fmt.Println("\nResolved names")
	for name, sym := range resolved {
		fmt.Printf(
			"USE %-10s %-12s at %d:%d  →  DEF %-10s %-12s at %d:%d\n",
			name.ID,
			sym.Kind,
			name.Pos.Start.Line,
			name.Pos.Start.Col,
			sym.Name,
			sym.Kind,
			sym.Span.Start.Line,
			sym.Span.Start.Col,
		)
	}
}

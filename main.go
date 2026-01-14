package main

import (
	"fmt"
	"os"
	"sort"

	"rahu/analyser"
	"rahu/parser"
	"rahu/utils"
)

func main() {
	input := utils.ParseFile("test.py")
	fmt.Printf("Going to parse input string\n%s\n\n\n", input)

	p := parser.New(input)
	module := p.Parse()

	fmt.Println("\nFirst pass, scopes:")
	global := analyser.BuildScopes(module)
	utils.PrintScopes(os.Stdout, global)

	fmt.Println("\nSecond pass, resolution:")
	errors, resolved := analyser.Resolve(module, global)

	if len(errors) > 0 {
		fmt.Printf("Found %d semantic error(s):\n", len(errors))
		for _, err := range errors {
			fmt.Printf("  Error at %v: %s\n", err.Span, err.Msg)
		}
	} else {
		fmt.Println("✓ No semantic errors")
	}

	fmt.Println("\nResolved names:")
	printResolvedNames(resolved)
}

func printResolvedNames(resolved map[*parser.Name]*analyser.Symbol) {
	// Convert to slice for sorting
	type resolution struct {
		name *parser.Name
		sym  *analyser.Symbol
	}

	var resolutions []resolution
	for name, sym := range resolved {
		resolutions = append(resolutions, resolution{name, sym})
	}

	// Sort by line then column for readable output
	sort.Slice(resolutions, func(i, j int) bool {
		if resolutions[i].name.Pos.Start.Line != resolutions[j].name.Pos.Start.Line {
			return resolutions[i].name.Pos.Start.Line < resolutions[j].name.Pos.Start.Line
		}
		return resolutions[i].name.Pos.Start.Col < resolutions[j].name.Pos.Start.Col
	})

	for _, r := range resolutions {
		name := r.name
		sym := r.sym

		// Determine scope description
		scopeDesc := "local"
		if sym.Span.Start.Line == 0 && sym.Span.Start.Col == 0 {
			scopeDesc = "builtin"
		}

		// Format output
		if scopeDesc == "builtin" {
			fmt.Printf(
				"  %d:%-3d  %-12s → builtin %s\n",
				name.Pos.Start.Line,
				name.Pos.Start.Col,
				name.ID,
				sym.Name,
			)
		} else {
			fmt.Printf(
				"  %d:%-3d  %-12s → %s (defined at %d:%d)\n",
				name.Pos.Start.Line,
				name.Pos.Start.Col,
				name.ID,
				sym.Name,
				sym.Span.Start.Line,
				sym.Span.Start.Col,
			)
		}
	}
}

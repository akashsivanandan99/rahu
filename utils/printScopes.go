package utils

import (
	"fmt"
	"io"
	"strings"

	"rahu/analyser"
)

func PrintScopes(w io.Writer, scope *analyser.Scope) {
	printScope(w, scope, 0)
}

func printScope(w io.Writer, s *analyser.Scope, indent int) {
	prefix := strings.Repeat(" ", indent)

	fmt.Fprintf(w, "%sScope(%s)\n", prefix, s.Kind)
	for _, sym := range s.Symbols {
		fmt.Fprintf(w, "%s  %s : %s\n", prefix, sym.Name, sym.Kind)
	}

	for _, child := range s.Children {
		printScope(w, child, indent+2)
	}
}

package analyser

import (
	"testing"

	"rahu/parser"
)

func TestFunctionScope(t *testing.T) {
	input := `
        def f(a, b):
            x = 1
    `

	p := parser.New(input)
	module := p.Parse()

	global := BuildScopes(module)

	if _, ok := global.Symbols["f"]; !ok {
		t.Fatal("missing function symbol f")
	}

	fn := global.Children[0]

	for _, name := range []string{"a", "b", "x"} {
		if _, ok := fn.Symbols[name]; !ok {
			t.Fatalf("missing symbol %s", name)
		}
	}
}

func TestSimpleResolution(t *testing.T) {
	src := `
	x = 1
	y = x
	`

	p := parser.New(src)
	module := p.Parse()
	global := BuildScopes(module)

	errs, resolved := Resolve(module, global)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %+v", errs)
	}

	var xUses []*parser.Name
	collectNames(module, &xUses)

	if resolved[xUses[0]].Name != "x" {
		t.Fatal("x did not resolve to symbol x")
	}
}

func collectNames(n parser.Node, out *[]*parser.Name) {
	switch v := n.(type) {

	// ----- statements -----

	case *parser.Module:
		for _, s := range v.Body {
			collectNames(s, out)
		}

	case *parser.Assign:
		for _, t := range v.Targets {
			collectNames(t, out)
		}
		collectNames(v.Value, out)

	case *parser.FunctionDef:
		for _, arg := range v.Args {
			if arg.Default != nil {
				collectNames(arg.Default, out)
			}
		}
		for _, s := range v.Body {
			collectNames(s, out)
		}

	case *parser.ExprStmt:
		collectNames(v.Value, out)

	case *parser.If:
		collectNames(v.Test, out)
		for _, s := range v.Body {
			collectNames(s, out)
		}
		for _, s := range v.Orelse {
			collectNames(s, out)
		}

	case *parser.For:
		collectNames(v.Target, out)
		collectNames(v.Iter, out)
		for _, s := range v.Body {
			collectNames(s, out)
		}

	case *parser.WhileLoop:
		collectNames(v.Test, out)
		for _, s := range v.Body {
			collectNames(s, out)
		}

	case *parser.Return:
		if v.Value != nil {
			collectNames(v.Value, out)
		}

	// ----- expressions -----

	case *parser.Name:
		*out = append(*out, v)

	case *parser.BinOp:
		collectNames(v.Left, out)
		collectNames(v.Right, out)

	case *parser.UnaryOp:
		collectNames(v.Operand, out)

	case *parser.Call:
		collectNames(v.Func, out)
		for _, a := range v.Args {
			collectNames(a, out)
		}

	case *parser.Compare:
		collectNames(v.Left, out)
		for _, r := range v.Right {
			collectNames(r, out)
		}

	case *parser.Tuple:
		for _, e := range v.Elts {
			collectNames(e, out)
		}

	case *parser.List:
		for _, e := range v.Elts {
			collectNames(e, out)
		}

	case *parser.BooleanOp:
		for _, e := range v.Values {
			collectNames(e, out)
		}

	// literals: nothing to do
	case *parser.Number, *parser.String, *parser.Boolean:
		return
	}
}

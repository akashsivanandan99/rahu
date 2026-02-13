package server

import (
	"rahu/lsp"
	"rahu/parser"
)

func ToRange(r parser.Range) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      r.Start.Line - 1,
			Character: r.Start.Col - 1,
		},
		End: lsp.Position{
			Line:      r.End.Line - 1,
			Character: r.End.Col - 1,
		},
	}
}

func FromRange(r lsp.Range) parser.Range {
	return parser.Range{
		Start: parser.Position{
			Line: r.Start.Line,
			Col:  r.Start.Character,
		},
		End: parser.Position{
			Line: r.End.Line + 1,
			Col:  r.End.Character + 1,
		},
	}
}

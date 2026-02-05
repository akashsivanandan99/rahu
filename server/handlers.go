package server

import (
	"rahu/jsonrpc"

	"rahu/lsp"
)

func (s *Server) DidOpen(p *lsp.DidOpenTextDocumentParams) {
	s.Open(p.TextDocument)
}

func (s *Server) Hover(p *lsp.HoverParams) (*lsp.Hover, *jsonrpc.Error) {
	doc := s.Get(p.TextDocument.URI)
	if doc == nil {
		return nil, jsonrpc.InvalidParamsError(nil)
	}

	return &lsp.Hover{
		Contents: "ok",
	}, nil
}

func (s *Server) DidChange(p *lsp.DidChangeTextDocumentParams) {
	s.ApplyFullChange(
		p.TextDocument.URI,
		p.ContentChange,
		p.TextDocument.Version,
	)
}

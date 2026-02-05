package server

import (
	"sync"

	"rahu/lsp"
)

type Server struct {
	mu           sync.RWMutex
	docs         map[lsp.DocumentURI]*Document
	capabilities lsp.ClientCapabilities
}

func New() *Server {
	return &Server{
		docs: make(map[lsp.DocumentURI]*Document),
	}
}

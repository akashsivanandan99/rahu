package server

import (
	"rahu/jsonrpc"
)

func Register(s *Server) {
	jsonrpc.RegisterRequest(
		"initialize",
		jsonrpc.AdaptRequest(s.Initialize),
	)
	jsonrpc.RegisterNofication(
		"textDocument/didOpen",
		jsonrpc.AdaptNotification(s.DidOpen),
	)

	jsonrpc.RegisterRequest(
		"textDocument/hover",
		jsonrpc.AdaptRequest(s.Hover),
	)

	jsonrpc.RegisterNofication(
		"textDocument/didChange",
		jsonrpc.AdaptNotification(s.DidChange),
	)
}

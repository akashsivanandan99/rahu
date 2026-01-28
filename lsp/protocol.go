// Package lsp defines Go representations of the Language Server Protocol (LSP)
// data structures for JSON-RPC communication. The types in this package are
// designed to closely follow the LSP specification and preserve wire-level
// compatibility.
package lsp

import (
	j "encoding/json"
)

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type DocumentURI string

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

type VersionedDocumentIdentifier struct {
	URI     DocumentURI `json:"uri"`
	Version int         `json:"version"`
}

type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type DocumentFilter struct {
	Language string `json:"language"`
	Scheme   string `json:"scheme"`
	Pattern  string `json:"pattern"`
}

type DocumentSelector []DocumentFilter

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type OptionalVersionedTextDocumentIdentifier struct {
	URI     DocumentURI `json:"uri"`
	Version *int        `json:"version,omitempty"`
}

type ChangeAnnotation struct {
	Label             string `json:"label"`
	NeedsConfirmation bool   `json:"needsConfirmation,omitempty"`
	Description       string `json:"description,omitempty"`
}

type ChangeAnnotationIdentifier string

type AnnotatedTextEdit struct {
	Range        Range                      `json:"range"`
	NewText      string                     `json:"newText"`
	AnnotationID ChangeAnnotationIdentifier `json:"annotationId"`
}

type TextDocumentEdit struct {
	TextDocument OptionalVersionedTextDocumentIdentifier `json:"textDocument"`
	Edits        j.RawMessage                            `json:"edits"`
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

type WorkspaceEdit struct {
	Chagnes          map[DocumentURI][]TextEdit                      `json:"changes,omitempty"`
	DocumentChanges  []TextDocumentEdit                              `json:"documentChanges,omitempty"`
	ChangeAnnotation map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotaion,omitempty"`
}

type Diagnostic struct {
	Range    Range  `json:"Range"`
	Severity int    `json:"severity,omitempty"`
	Code     any    `json:"code,omitempty"`
	Source   string `json:"source,omitempty"`
	Message  string `json:"message"`
}

type LocationLink struct {
	OriginSelectionRange *Range      `json:"originSelectionRange,omitempty"`
	TargetURI            DocumentURI `json:"TargetUri"`
	TargetRange          Range       `json:"TargetRange"`
	TargetSelectionRange Range       `json:"TargetSelectionRange"`
}

type ProgressToken j.RawMessage

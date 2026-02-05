// Package lsp defines Go representations of the Language Server Protocol (LSP)
// data structures for JSON-RPC communication. The types in this package are
// designed to closely follow the LSP specification and preserve wire-level
// compatibility.
//
// This package provides types for positions, ranges, documents, diagnostics,
// and other LSP constructs, ensuring seamless integration with LSP clients.
package lsp

import (
	j "encoding/json"
)

type Position struct {
	// Line represents the zero-based line number in a document.
	Line int `json:"line"`

	// Character represents the zero-based character offset within the line.
	Character int `json:"character"`
}

type DocumentURI string // DocumentURI is a string representing the URI of a document.

type Range struct {
	// Start represents the starting position of the range.
	Start Position `json:"start"`

	// End represents the ending position of the range.
	End Position `json:"end"`
}

type TextDocumentIdentifier struct {
	// URI is the unique identifier of the text document.
	URI DocumentURI `json:"uri"`
}

type VersionedDocumentIdentifier struct {
	// URI is the unique identifier of the text document.
	URI DocumentURI `json:"uri"`

	// Version is the version number of the document.
	Version int `json:"version"`
}

type TextDocumentItem struct {
	// URI is the unique identifier of the text document.
	URI DocumentURI `json:"uri"`

	// LanguageID specifies the language of the document (e.g., "python").
	LanguageID string `json:"languageId"`

	// Version is the version number of the document.
	Version int `json:"version"`

	// Text contains the full content of the document.
	Text string `json:"text"`
}

type TextDocumentPositionParams struct {
	// TextDocument identifies the document in which the position is specified.
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	// Position specifies the location within the document.
	Position Position `json:"position"`
}

type DocumentFilter struct {
	// Language specifies the language of the document (e.g., "go").
	Language string `json:"language"`

	// Scheme specifies the URI scheme (e.g., "file").
	Scheme string `json:"scheme"`

	// Pattern specifies the glob pattern to match document paths.
	Pattern string `json:"pattern"`
}

type DocumentSelector []DocumentFilter // DocumentSelector is a list of filters used to match documents.

type TextEdit struct {
	// Range specifies the range of text to be replaced.
	Range Range `json:"range"`

	// NewText is the text to replace the specified range with.
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
	Edits        any                                     `json:"edits"` // []TextEdit | []AnnotatedTextEdit
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

type WorkspaceEdit struct {
	Changes           map[DocumentURI][]TextEdit                      `json:"changes,omitempty"`
	DocumentChanges   []TextDocumentEdit                              `json:"documentChanges,omitempty"`
	ChangeAnnotations map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotations,omitempty"`
}

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityInformation
	SeverityHint
)

type Diagnostic struct {
	Range    Range    `json:"range"`
	Severity Severity `json:"severity,omitempty"`
	Code     any      `json:"code,omitempty"`
	Source   string   `json:"source,omitempty"`
	Message  string   `json:"message"`
}

type LocationLink struct {
	OriginSelectionRange *Range      `json:"originSelectionRange,omitempty"`
	TargetURI            DocumentURI `json:"targetUri"`
	TargetRange          Range       `json:"targetRange"`
	TargetSelectionRange Range       `json:"targetSelectionRange"`
}

type ProgressToken any

type ProgressParams struct {
	Token ProgressToken `json:"token"`
	Value any           `json:"value"`
}

type HoverParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type MarkupKind string

const (
	MarkupPlaintext MarkupKind = "plaintext"
	MarkupMarkdown  MarkupKind = "markdown"
)

type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

type MarkedString any

type Hover struct {
	Contents any    `json:"contents"`
	Range    *Range `json:"range,omitempty"`
}

type WorkDoneProgressParams struct {
	WorkDoneToken ProgressToken `json:"workDoneToken,omitempty"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type ResourceOperationKind string

const (
	CREATE ResourceOperationKind = "create"
	RENAME ResourceOperationKind = "rename"
	DELETE ResourceOperationKind = "delete"
)

type FailureHandlingKind string

type ChangeAnnotationSupport struct {
	GroupsOnLabel *bool `json:"groupsOnLabel,omitempty"`
}

const (
	ABORT                 FailureHandlingKind = "abort"
	TRANSACTIONAL         FailureHandlingKind = "transactional"
	TEXTONLYTRANSACTIONAL FailureHandlingKind = "textOnlyTransactional"
	UNDO                  FailureHandlingKind = "undo"
)

type WorkspaceEditClientCapabilities struct {
	DocumentChanges         *bool                    `json:"documentChanges,omitempty"`
	ResourceOperations      *[]ResourceOperationKind `json:"resourceOperations,omitempty"`
	FailureHandling         *FailureHandlingKind     `json:"failureHandlingKind,omitempty"`
	NormalizesLineEndings   *bool                    `json:"normalizesLineEndings,omitempty"`
	ChangeAnnotationSupport *ChangeAnnotationSupport `json:"changeAnnotationSupport,omitempty"`
}

type DidChangeConfigurationParams struct {
	Settings any `json:"settings"`
}

type DidChangeConfigurationClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
}

type DidChangeWatchedFilesClientCapabilities struct {
	DynamicRegistration    *bool `json:"dynamicRegistration,omitempty"`
	RelativePatternSupport *bool `json:"relativePatternSupport,omitempty"`
}

type WatchKind int

const (
	Create WatchKind = 1
	Change WatchKind = 2
	Delete WatchKind = 4
)

type WorkspaceFolder struct {
	URI  DocumentURI `json:"uri"`
	Name string      `json:"name"`
}

type Pattern string

type RelativePattern struct {
	BaseURI any     `json:"baseUri"`
	Pattern Pattern `json:"pattern"`
}

type GlobPattern any

type FileSystemWatcher struct {
	Kind        *WatchKind  `json:"kind,omitempty"`
	GlobPattern GlobPattern `json:"globPattern"`
}

type DidChangeWatchedFilesRegistrationOptions struct {
	Watchers []FileSystemWatcher `json:"watchers"`
}

type WorkspaceClientCapabilities struct {
	ApplyEdit              *bool                                     `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditClientCapabilities          `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`
	DidChangeWatchedFiles  *DidChangeWatchedFilesClientCapabilities  `json:"didChangeWatchedFiles,omitempty"`
}

type TextDocumentClientCapabilities struct{}

type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Window       j.RawMessage                    `json:"window,omitempty"`
	General      j.RawMessage                    `json:"general,omitempty"`
}

type InitializeParams struct {
	WorkDoneToken         ProgressToken      `json:"workDoneToken,omitempty"`
	ProcessID             *int               `json:"processId,omitempty"`
	ClientInfo            *ClientInfo        `json:"clientInfo,omitempty"`
	Locale                *string            `json:"locale,omitempty"`
	RootPath              *string            `json:"rootPath,omitempty"`
	RootURI               *DocumentURI       `json:"rootUri,omitempty"`
	InitializationOptions j.RawMessage       `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ClientInfo        `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncKind `json:"textDocumentSync"`
	HoverProvider      bool                 `json:"hoverProvider"`
	CompletionProvider map[string]any       `json:"completionProvider,omitempty"`
}

type TextDocumentSyncKind int

const (
	TDSKNone        TextDocumentSyncKind = 0
	TDSKFull        TextDocumentSyncKind = 1
	TDSKIncremental TextDocumentSyncKind = 2
)

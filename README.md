# rahu
A Python Language Server Protocol (LSP) implementation written in Go.

## About
Rahu is a learning project to understand how LSPs work by building one from scratch. It implements lexical analysis, parsing, and language server features for Python code.

## Current Status
🚧 **Work in Progress** 🚧

### Completed
- ✅ **Lexer**: Full Python tokenization with INDENT/DEDENT support
  - All Python operators (single and multi-character)
  - String literals (single-line and multi-line)
  - Number and identifier recognition
  - Proper indentation tracking with tab/space consistency checking
  - Basic test coverage

### In Progress
- 🔄 **Parser**: Abstract Syntax Tree (AST) generation
- 🔄 **Language Server**: LSP protocol implementation

### Planned
- ⏳ Semantic analysis
- ⏳ Go-to-definition
- ⏳ Auto-completion
- ⏳ Diagnostics

## Project Structure
```
rahu/
├── lexer/          # Tokenization
├── parser/         # AST generation (coming soon)
├── server/         # LSP server (coming soon)
└── utils/          # Utilities
```

## Building & Running

### Prerequisites
- Go 1.21 or higher

### Build
```bash
go build
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run lexer tests with verbose output
go test ./lexer -v
```

### Try the Lexer
```bash
go run main.go test.py
```

## Example

Given this Python code:
```python
def add(x, y):
    return x + y
```

The lexer produces:
```
NAME         "def"
NAME         "add"
LPAR         "("
NAME         "x"
COMMA        ","
NAME         "y"
RPAR         ")"
COLON        ":"
NEWLINE      "\n"
INDENT       ""
NAME         "return"
NAME         "x"
PLUS         "+"
NAME         "y"
NEWLINE      "\n"
DEDENT       ""
EOF          ""
```

## Development

This project follows a learning-focused approach:
1. Lexer (tokenization) ✅
2. Parser (AST generation) 🔄
3. Semantic analysis
4. LSP server implementation

Each phase builds on the previous one.

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [Python Grammar](https://docs.python.org/3/reference/grammar.html)
- [Python Token Documentation](https://docs.python.org/3/library/token.html)

## License
MIT

## Author
Akash Sivanandan

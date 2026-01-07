# rahu

A Python Language Server Protocol (LSP) implementation written in Go.

## About

Rahu is a learning project to understand how LSPs work by building one from scratch. It implements lexical analysis, parsing, and language server features for Python code.

## Current Status

🚧 **Work in Progress** 🚧

### Completed

- ✅ **Lexer**: Full Python tokenization with INDENT/DEDENT support
  - All Python operators (single and multi-character: `+`, `-`, `*`, `/`, `//`, `%`, `==`, `!=`, `<`, `>`, `<=`, `>=`, etc.)
  - String literals (single-line and multi-line with `"""` and `'''`)
  - Number, identifier, and keyword recognition
  - Proper indentation tracking with tab/space consistency checking
  - Position tracking (line and column numbers)
  - Comprehensive test coverage

- ✅ **Parser**: Complete recursive descent parser with AST generation
  - **Statements**: assignments, if/elif/else, for loops, while loops, function definitions, return, break, continue
  - **Expressions**: binary operations, comparisons, boolean operations (and/or), unary operations, function calls, lists, tuples
  - **Operator precedence**: Correctly handles arithmetic, comparison, and boolean operator precedence
  - **Advanced features**: 
    - Function default arguments
    - Comparison chaining (`1 < x < 10`)
    - Tuple unpacking in for loops (`for x, y in pairs`)
    - Nested control structures
  - 40+ test cases covering all major features

### In Progress

- 🔄 **Semantic Analysis**: Symbol tables, scope resolution, type checking

### Planned

- ⏳ **Language Server**: LSP protocol implementation
- ⏳ Go-to-definition
- ⏳ Auto-completion
- ⏳ Diagnostics and error reporting
- ⏳ Hover information

## Project Structure

```
rahu/
├── lexer/          # Tokenization
│   ├── lexer.go    # Lexer implementation
│   ├── token.go    # Token definitions
│   ├── maps.go     # Operator and keyword maps
│   └── lexer_test.go
├── parser/         # AST generation
│   ├── parser.go   # Parser implementation
│   ├── ast.go      # AST node definitions
│   └── parser_test.go
├── utils/          # Helper utilities
│   └── printAST.go # AST visualization
└── main.go         # Entry point
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

# Run parser tests with verbose output
go test ./parser -v
```

### Try the Parser

```bash
go run main.go
```

## Example

Given this Python code:

```python
def fibonacci(n):
    if n <= 1:
        return n
    return fibonacci(n - 1) + fibonacci(n - 2)

result = fibonacci(10)
```

The lexer tokenizes it into:

```
DEF          "def"
NAME         "fibonacci"
LPAR         "("
NAME         "n"
RPAR         ")"
COLON        ":"
NEWLINE      "\n"
INDENT       ""
IF           "if"
NAME         "n"
LESSEQUAL    "<="
NUMBER       "1"
...
```

And the parser generates an AST:

```
Module:
  FunctionDef:
    Name: fibonacci
    Args: [n]
    Body:
      If:
        Test: Compare(n <= 1)
        Body:
          Return: Name(n)
        Orelse: []
      Return:
        BinOp:
          Left: Call(fibonacci, [n - 1])
          Op: +
          Right: Call(fibonacci, [n - 2])
  Assign:
    Target: Name(result)
    Value: Call(fibonacci, [10])
```

## Features Supported

### Statements
- Variable assignments: `x = 5`
- Function definitions: `def foo(x, y=10):`
- Control flow: `if`/`elif`/`else`, `for`, `while`
- Loop control: `break`, `continue`
- Return statements: `return x + y`

### Expressions
- Arithmetic: `+`, `-`, `*`, `/`, `//`, `%`
- Comparisons: `==`, `!=`, `<`, `>`, `<=`, `>=` (with chaining: `1 < x < 10`)
- Boolean: `and`, `or`, `not`
- Unary: `-x`, `+x`, `not x`
- Function calls: `foo(1, 2, 3)`
- Lists: `[1, 2, 3]`, `[[1, 2], [3, 4]]`
- Parentheses: `(a + b) * c`

### Advanced Features
- **Operator precedence**: Correctly parses `1 + 2 * 3` as `1 + (2 * 3)`
- **Default arguments**: `def greet(name, greeting="Hello")`
- **Tuple unpacking**: `for x, y in pairs:`
- **Nested structures**: Complex nested if/for/while statements

## Development

This project follows a learning-focused approach:

1. **Lexer** (tokenization) ✅ Complete
2. **Parser** (AST generation) ✅ Complete
3. **Semantic analysis** 🔄 In Progress
4. **LSP server implementation** ⏳ Planned

Each phase builds on the previous one, with a focus on understanding the fundamentals of language implementation.

## Technical Highlights

- **Hand-written lexer**: No lexer generators, built from scratch for learning
- **Recursive descent parser**: Clean, understandable parsing logic
- **Pratt parsing**: Efficient operator precedence handling
- **Indentation-aware**: Proper Python-style indentation handling with INDENT/DEDENT tokens
- **Comprehensive tests**: 40+ test cases ensuring correctness

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [Python Grammar](https://docs.python.org/3/reference/grammar.html)
- [Python AST Documentation](https://docs.python.org/3/library/ast.html)
- [Crafting Interpreters](https://craftinginterpreters.com/)

## License

MIT

## Author

Akash Sivanandan

# expr-eval

Mathematical and multi-type expression evaluator with bytecode compilation (CLI).

Parses expressions with arithmetic, comparison, and logical operators, compiles to stack-machine bytecode, and evaluates with variable bindings. Supports numeric, boolean, and string types with a rich built-in function library.

## Architecture

```
lexer → parser → AST
                   ├── eval (tree-walk interpreter)
                   ├── compile → optimize → vm (bytecode execution)
                   ├── validate (static analysis)
                   └── format (pretty-print / RPN)
         env (scoped variables + constants)
         types (multi-type value system)
         builtin (math + string + logic function library)
```

Packages:

| Package | Role |
|---------|------|
| `internal/lexer` | Tokenization: numbers, identifiers, operators, punctuation |
| `internal/parser` | Recursive descent parser producing AST nodes |
| `internal/ast` | Extended AST node types with tree traversal utilities |
| `internal/eval` | Tree-walking evaluator (legacy, float64-only) |
| `internal/types` | Multi-type value system (number/bool/string/null) with coercion |
| `internal/builtin` | Function registry: 22 math, 13 string, 12 logic functions |
| `internal/compile` | Bytecode compiler with constant folding and strength reduction |
| `internal/vm` | Stack virtual machine for compiled bytecode execution |
| `internal/validate` | Static analysis: undefined vars, arity checks, complexity |
| `internal/format` | AST pretty-printing, parenthesization, RPN output |
| `internal/env` | Scoped variable environment with built-in constants |

## Usage

```bash
expr-eval "2 * pi * r" -r 5
expr-eval "(a + b) ^ 2 / sqrt(c)" -a 3 -b 4 -c 25
expr-eval "min(x, max(y, 0))" -x=10 -y=-3
```

## Supported Features

- Arithmetic: `+`, `-`, `*`, `/`, `%`, `^`
- Functions: `abs`, `sqrt`, `ceil`, `floor`, `round`, `log`, `exp`, `sin`, `cos`, `tan`, `min`, `max`, `clamp`, `pow`, `hypot`, etc.
- Variables with scoped environment and built-in constants (`pi`, `e`, `tau`)
- Bytecode compilation with peephole optimization (constant folding, dead code elimination)
- Static validation (undefined variable detection, function arity checking)

## Build & Test

```bash
export GOTOOLCHAIN=local CGO_ENABLED=0
go build ./...
go test ./...
```

## License

MIT

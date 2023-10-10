# DEVELOP

## Description

This package is designed to traverse an JSON Abstract Syntax Tree (AST) and convert its interpretation into various output formats.

## Constructing the AST

### Idea

The core idea of this package is to use a JSON AST as an intermediary layer. This design enables JSON to drive data conversion, facilitating conversions between any input to JSON or from JSON to any desired output. As the AST is extended in terms of syntax or capabilities, it could pave the way for a robust templating engine.

Each AST node has a `Visit` method, allowing the process of consuming the AST to be delegated to the visitor's implementation.

```golang
type JsonNode interface {
    GetNodeType() AST_NODETYPE
    Visit(visitor JsonVisitor) error
    String() string
    ShouldOmitEmpty() bool
}
```

### Construction procedure

The `constructor` package dictates the procedure for building the AST based on different input sources.

Key interfaces are defined in the `protocol.go` file:

```golang
type TokenProvider interface {
    ReadBool() (bool, error)
    ReadNull() error
    ReadNumber() (float64, error)
    ReadString() ([]byte, error)
    ReadVariable() ([]byte, error)
    GetNextTokenType() (token.TokenType, error)
}

type ASTManager interface {
    RecordStateValue(valueType ast.AST_NODETYPE, nodeValue interface{}) error
    GetAST() ast.JsonNode
    HasComplete() bool
    TopElementType() (ast.AST_NODETYPE, error)
    HasOpenElements() bool
}

type ASTBuilder interface {
    ASTManager
    TokenProvider
}
```

The `ASTBuilder` consists of two parts: `TokenProvider` for tokenizing input and `ASTManager` for recording the corresponding values of tokens.

### Symbol-Driven Construction

While AST creation is symbol-driven, this detail is abstracted away from the tokenizer and managed internally within the `ASTBuilder`.

The `golang` implementation, for instance, is also symbol-driven. When constructing the AST, `GetNextTokenType` is the tokenizer's core function to consume the input, when a json symbol comes, it will construct the AST internally The `RecordStateValue` is then invoked by the tokenizer to insert the value at the appropriate location in the AST.

```golang
func (t *astGolangConstructor) GetNextTokenType() (token.TokenType, error) {

    nextTokenType, err := t.provider.GetNextTokenType()
    if err != nil {
        return token.TOKEN_DUMMY, err
    }

    if token.IsSymbolToken(nextTokenType) {
        err = t.astConstructor.RecordSyntaxSymbol(nextTokenType)
        if err != nil {
            return token.TOKEN_DUMMY, err
        }
    }

    return nextTokenType, nil
}

func (i *astGolangConstructor) RecordSyntaxSymbol(b token.TokenType) error {
    switch b {
    case token.TOKEN_LEFT_BRACE:
        return i.ast.CreateNewASTNode(ast.AST_OBJECT, nil)
    case token.TOKEN_LEFT_BRACKET:
        return i.ast.CreateNewASTNode(ast.AST_ARRAY, nil)
    case token.TOKEN_RIGHT_BRACKET:
        fallthrough
    case token.TOKEN_RIGHT_BRACE:
        return i.ast.EncloseLatestElements()
    default:
        return ErrorIncorrectSyntaxSymbolForConstructAST
    }
    return nil
}
```

### KV Pair Nodes

The design of the KV pair nodes reflects the key-value nature of JSON. During object construction, the AST will manage a `kv pair` node, encapsulating the primitive value node. This design simplifies stack management, ensuring only the `object`, `array`, and `kv pair` nodes remain on the stack.

## Marshaling

When marshaling Go data, the tokenizer is used to process the Go struct and convert it to the JSON AST.

Once the AST is constructed, it can be traversed and converted into any compatible data output. Implementing the marshaler in a visitor mode streamlines this process:

```golang
func InterpretAST(node ast.JsonNode, variables map[string]interface{}, marshaler ast.MarshalerFunc) ([]byte, error) {

    visitor := NewASTInterpreter(variables, marshaler)
    visitor.stackNode.Push(node)

    for {
        node, err := visitor.stackNode.Pop()
        if err != nil {
            if err != util.ErrorEndOfStack {
                return nil, err
            }
            break
        }

        err = node.Visit(visitor)
        if err != nil {
            return nil, err
        }
    }

    if visitor.getSymbolLength() > 0 {
        return nil, ErrorInterpreSymbolFailure
    }
    return visitor.GetOutput(), nil
}
```

## Unmarshaling

Unmarshaling involves converting byte input to a Go struct. This process mirrors marshaling, as both involve consuming the AST. The primary difference is

 the input source; while marshaling uses Go data to create the AST, unmarshaling uses byte data.

When unmarshaling, the AST is converted to the designated object, similarly to marshaling.


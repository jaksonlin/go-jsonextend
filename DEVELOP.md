# DEVELOP

## description

this package is design to traverse the AST and interpret it into different type of output

## construction of AST

### idea

the design of this package is to use a JSON AST to traverse as a middle-layer, this can make json as a core to drive data conversion; being able to convert any input into json; or convert json into any kind of output. this is the main idea! when later we extend the syntax or capability of the AST, we will be able to create a powerful templating engine.

every AST node have a visitor function `Visit`, so that we can proxy the procedure of consuming the AST information to the implementation of visitor. 

```golang
type JsonNode interface {
    GetNodeType() AST_NODETYPE
    Visit(visitor JsonVisitor) error
    String() string
    ShouldOmitEmpty() bool
}
```

when the visiting of AST node is implemented, there's only few work left to contruct the AST: symbol driven creation.

### constructing

the package name `constructor` means the construction procedure of AST based on different input source.

and we defines the key interface in `protocol.go` file

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

the ASTBuilder contains 2 parts, `TokenProvider` interface for tokenize the input, and `ASTManager` to record the token's cooresponding value.

while one will need to implement all the function for `TokenProvider`, the `ASTManager` the main part is `RecordStateValue` others are ast proxy function for tokenizer to do state check/change.

### symbol driven construction

One important things to note is that, while the AST creation is symbol driven, this fact is not reveal to the tokenizer, this is being kept within the ASTBuilder, (one can do it in a none-symbol driven approach). 

currently the `golang` impllementation is also symbol driven, because it is straitforward that a struct/map is symbol `{` and array/slice is symbol `[`. the difficult part will rely on when to signal a closing symbol (`}` or `]`), this is left for implemenation of how one will consume the input.

for exmaple, currently the package use a stack to traverse the AST(none-recrusive), and for each processing of object/array, we will push the closing symbol prior to pushing any value entity into the stack, so that when the object finishes its value consuming, the tokenizer will receive a closing symbol.

and here is a draft of creating AST to start with, the `GetNextTokenType` is tokenizer's key function to consume the input, when a json symbol comes, we will construct the AST internally, and the `RecordStateValue` will then being called by tokenizer to inject the value into the correct location we create in `RecordSyntaxSymbol` when the token provider first sees the symbol ;-)

this part is the same for golang/bytebase implementation.

``` golang
// put the store to syntax symbol here, to decouple the relation of reader and writer
func (t *astGolangConstructor) GetNextTokenType() (token.TokenType, error) {

    nextTokenType, err := t.provider.GetNextTokenType()
    if err != nil {
        return token.TOKEN_DUMMY, err
    }

    if token.IsSymbolToken(nextTokenType) { // note symbol token will be parse in the corresponding primitive value state
        err = t.astConstructor.RecordSyntaxSymbol(nextTokenType)
        if err != nil {
            return token.TOKEN_DUMMY, err
        }
    }

    return nextTokenType, nil
}
func (i *astGolangConstructor) RecordSyntaxSymbol(b token.TokenType) error {
    //routing base on symbol
    switch b {
    case token.TOKEN_LEFT_BRACE:
        return i.ast.CreateNewASTNode(ast.AST_OBJECT, nil)
    case token.TOKEN_LEFT_BRACKET:
        return i.ast.CreateNewASTNode(ast.AST_ARRAY, nil)
    case token.TOKEN_RIGHT_BRACKET:
        fallthrough
    case token.TOKEN_RIGHT_BRACE:
        err := i.ast.EncloseLatestElements()
        if err != nil {
            return err
        }
    default:
        return ErrorIncorrectSyntaxSymbolForConstructAST
    }
    return nil
}
```

### kv pair node

the kv pair node is design to reflect the fact of json (template) that the values are kv pair. during construction of an object, the ast stack state will hold an `object`->`kv pair`, it is designed in this way, so that when a primitive value node comes, it will be held inside the `kv pair node`, not on the stack, this can greatly ease the stack management: no primitive value node on the stack, only: `object`, `array`, `kvpair`. once the `kv pair`'s key and values are resolved, it will be push into the owner object node, therefore for an object having many kv pair, there's only 1 kv pair node for this object on the stack. ;-)

## marshaling

when marshaling a go data, we will use the tokenizer to tokenize the go struct and convert to Json AST as describe above. 

when the AST is done, we can traverse it and convert into any (json compatible) data output. this is a traditional tree traverse, the implementation use a none-recrusive approach. imeplement the marshaler in visitor mode can greatly ease the procedure:

```

func InterpretAST(node ast.JsonNode, variables map[string]interface{}, marshaler ast.MarshalerFunc) ([]byte, error) {
    // deep first traverse the AST

    visitor := NewASTInterpreter(variables, marshaler)
    visitor.stackNode.Push(node)

    for {
        node, err := visitor.stackNode.Pop()
        if err != nil {
            break
        }
        err = node.Visit(visitor)
        if err != nil {
            if err != util.ErrorEndOfStack {
                return nil, err
            } else {
                break
            }
        }

    }

    if visitor.getSymbolLength() > 0 {
        return nil, ErrorInterpreSymbolFailure
    }
    rs := visitor.GetOutput()
    return rs, nil
}
```

all you need to do is to check the ast node type and decide what to do with the values and output buffers/reflect.Value. ;-)

## unmarshaling

unmarshling the bytes input to go struct is basically the same as marshaling, as both of them are essential `consumption of AST`.

while the marshaling process use a go data as tokenizer's input to create the AST, the unmarshaling uses []byte. the rest are basically the same as [marshaling](#marshaling). when marshaling convert the AST to []byte, unmarshaling convert the AST to object you defines.


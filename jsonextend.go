package jsonextend

import (
	"io"

	"github.com/jaksonlin/go-jsonextend/interpreter"
)

// parse a jsonextend document with the variables into json bytes.
func Parse(reader io.Reader, variables map[string]interface{}) ([]byte, error) {
	return interpreter.ParseJsonExtendDocument(reader, variables)
}

// unmarshal a jsonextend document with the variables into a struct. should alied with json.Unmarshal
func Unmarshal(reader io.Reader, variables map[string]interface{}, out interface{}) error {
	return interpreter.Unmarshal(reader, variables, out)
}

// marshal a struct into json bytes. should alied with json.Marshal
func Marshal(v interface{}) ([]byte, error) {
	return interpreter.Marshal(v)
}

func MarshalWithVariables(v interface{}, variables map[string]interface{}) ([]byte, error) {
	return interpreter.MarshalWithVariables(v, variables)
}

func MarshalIntoTemplate(v interface{}) ([]byte, error) {
	return interpreter.MarshalIntoTemplate(v)
}

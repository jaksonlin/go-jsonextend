package jsonextend_test

import (
	"encoding/json"
)

var testMarshaler = func(v interface{}) ([]byte, error) { return json.Marshal(v) }

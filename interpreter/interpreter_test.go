package interpreter_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jaksonlin/go-jsonextend/interpreter"
	"github.com/jaksonlin/go-jsonextend/tokenizer"

	"net/http"
	_ "net/http/pprof"
)

func TestMain(m *testing.M) {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	time.Sleep(time.Second) // Give some time for the server to start

	os.Exit(m.Run())
}

func TestInterpreter(t *testing.T) {

	const sampleJson = `{
		"key1": "value1",
		"key2": 123,
		"key3": true,
		"key4": null,
		"key5": {
			"key6":  [
				"item1111",
				"item22222",
				"item33333"
			],
			"key7": 789,
			"key8": [
				"item1111333",
				"item22222333",
				"item33333333"
			],
			"key9": null,
			"key10": [
				"item1",
				"item2",
				"item3",
				${someVar1}
			]
		},
		"key11": [
			"item1",
			"item2",
			"item3",
			123,
			true,
			null,
			{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			}
		],
		"key16": {
			"some161":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			},
			"some162":123,
			"some163":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			},
			"some164":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			}
		},
		"key17": 123.4,
		"${key18}": "value18",
		"key19": "oh \"${value19}\"",
		"key20": ${value20}

	}`
	reader := strings.NewReader(sampleJson)
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)

	err := sm.ProcessData()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sm.GetASTBuilder().HasOpenElements() {
		t.FailNow()
	}
	a := sm.GetAST()

	variableConfig := make(map[string]interface{})
	variableConfig["someVar2"] = 123
	variableConfig["someVar1"] = map[string]interface{}{"hello": 123, "world": 223}
	variableConfig["key18"] = "key18+"
	variableConfig["value19"] = "somevalue"
	variableConfig["value20"] = []int{1, 2, 3}
	rs, err := interpreter.PrettyInterpret(a, variableConfig, interpreter.Marshal)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fmt.Println(string(rs))

	var someMap map[string]interface{}
	err = json.Unmarshal(rs, &someMap)
	if err != nil {
		t.FailNow()
	}
	fmt.Printf("%#v", someMap)
	if len(someMap) != 11 {
		t.FailNow()
	}
}

func TestOmit(t *testing.T) {
	type AnotherData struct {
		Name string
	}

	type Data struct {
		StringField    string            `json:"stringField,omitempty"`
		IntField       int               `json:"intField,omitempty"`
		BoolField      bool              `json:"boolField,omitempty"`
		SliceField     []string          `json:"sliceField,omitempty"`
		MapField       map[string]string `json:"mapField,omitempty"`
		PointerField   *string           `json:"pointerField,omitempty"`
		StructField    AnotherData       `json:"structField,omitempty"`
		InterfaceField interface{}       `json:"interfaceField,omitempty"`
	}

	d := Data{
		StringField: "",
		IntField:    0,
		BoolField:   false,
		SliceField:  nil,
		MapField:    nil,
		PointerField: func() *string {
			var s string
			return &s
		}(),
		StructField:    AnotherData{},
		InterfaceField: nil,
	}
	result, err := interpreter.Marshal(d)
	if err != nil {
		t.FailNow()
	}
	result2, err := json.Marshal(d)
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(result, result2) {
		t.FailNow()
	}

}

func TestOmit2(t *testing.T) {
	type AnotherData struct {
		Name string
	}

	type Data struct {
		StringField    *string            `json:"stringField,omitempty"`
		IntField       *int               `json:"intField,omitempty"`
		BoolField      *bool              `json:"boolField,omitempty"`
		SliceField     *[]string          `json:"sliceField,omitempty"`
		MapField       *map[string]string `json:"mapField,omitempty"`
		PointerField   *string            `json:"pointerField,omitempty"`
		StructField    *AnotherData       `json:"structField,omitempty"`
		InterfaceField *interface{}       `json:"interfaceField,omitempty"`
	}

	d := Data{}
	result, err := interpreter.Marshal(d)
	if err != nil {
		t.FailNow()
	}
	result2, err := json.Marshal(d)
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(result, result2) {
		t.FailNow()
	}

}

func TestInterpreterAST(t *testing.T) {

	const sampleJson = `{
		"key1": "value1",
		"key2": 123,
		"key3": true,
		"key4": null,
		"key5": {
			"key6":  [
				"item1111",
				"item22222",
				"item33333"
			],
			"key7": 789,
			"key8": [
				"item1111333",
				"item22222333",
				"item33333333"
			],
			"key9": null,
			"key10": [
				"item1",
				"item2",
				"item3",
				${someVar1}
			]
		},
		"key11": [
			"item1",
			"item2",
			"item3",
			123,
			true,
			null,
			{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			}
		],
		"key16": {
			"some161":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			},
			"some162":123,
			"some163":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			},
			"some164":{
				"key12": "value12",
				"key13": 456,
				"key14": true,
				"key15": null,
				"someVar2": ${someVar2}
			}
		},
		"key17": 123.4,
		"${key18}": "value18",
		"key19": "oh \"${value19}\"",
		"key20": ${value20}

	}`
	reader := strings.NewReader(sampleJson)
	sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)

	err := sm.ProcessData()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sm.GetASTBuilder().HasOpenElements() {
		t.FailNow()
	}
	a := sm.GetAST()

	variableConfig := make(map[string]interface{})
	variableConfig["someVar2"] = 123
	variableConfig["someVar1"] = map[string]interface{}{"hello": 123, "world": 223}
	variableConfig["key18"] = "key18+"
	variableConfig["value19"] = "somevalue"
	variableConfig["value20"] = []int{1, 2, 3}
	rs, err := interpreter.InterpretAST(a, variableConfig, testMarshaler)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fmt.Println(rs)

	var someMap map[string]interface{}
	err = json.Unmarshal(rs, &someMap)
	if err != nil {
		t.FailNow()
	}
	fmt.Printf("%#v", someMap)
	if len(someMap) != 11 {
		t.FailNow()
	}

	if someMap["key1"] != "value1" {
		t.FailNow()
	}
	if someMap["key2"] != 123.0 {
		t.FailNow()
	}
	if someMap["key3"] != true {
		t.FailNow()
	}
	if someMap["key4"] != nil {
		t.FailNow()
	}
	if someMap["key5"] == nil {
		t.FailNow()
	}
	key5Item := someMap["key5"].(map[string]interface{})
	if key5Item["key6"] == nil || len(key5Item["key6"].([]interface{})) != 3 {
		t.FailNow()
	}

	if key5Item["key7"] != 789.0 {
		t.FailNow()
	}

	if key5Item["key8"] == nil || len(key5Item["key8"].([]interface{})) != 3 {
		t.FailNow()
	}

	if key5Item["key9"] != nil {
		t.FailNow()
	}

	if key5Item["key10"] == nil || len(key5Item["key10"].([]interface{})) != 4 {
		t.FailNow()
	}

	if someMap["key11"] == nil {
		t.FailNow()
	}

	key11Item := someMap["key11"].([]interface{})
	if len(key11Item) != 7 {
		t.FailNow()
	}
	if key11Item[0] != "item1" {
		t.FailNow()
	}
	if key11Item[1] != "item2" {
		t.FailNow()
	}
	if key11Item[2] != "item3" {
		t.FailNow()
	}
	if key11Item[3] != 123.0 {
		t.FailNow()
	}
	if key11Item[4] != true {
		t.FailNow()
	}
	if key11Item[5] != nil {
		t.FailNow()
	}
	if key11Item[6] == nil {
		t.FailNow()
	}
	key11Item6 := key11Item[6].(map[string]interface{})
	if key11Item6["key12"] != "value12" {
		t.FailNow()
	}
	if key11Item6["key13"] != 456.0 {
		t.FailNow()
	}
	if key11Item6["key14"] != true {
		t.FailNow()
	}
	if key11Item6["key15"] != nil {
		t.FailNow()
	}
	if key11Item6["someVar2"] != 123.0 {
		t.FailNow()
	}

	if someMap["key16"] == nil {
		t.FailNow()
	}

	key16Item := someMap["key16"].(map[string]interface{})
	if key16Item["some161"] == nil {
		t.FailNow()
	}

	if key16Item["some162"] != 123.0 {
		t.FailNow()
	}

	if key16Item["some163"] == nil {
		t.FailNow()
	}

	if key16Item["some164"] == nil {
		t.FailNow()
	}

	if someMap["key17"] != 123.4 {
		t.FailNow()
	}
	if someMap["key18+"] != "value18" {
		t.FailNow()
	}
	if someMap["key19"] != "oh \"somevalue\"" {
		t.FailNow()
	}
	if someMap["key20"] == nil {
		t.FailNow()
	}
}

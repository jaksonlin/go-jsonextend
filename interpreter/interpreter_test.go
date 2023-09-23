package interpreter_test

import (
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
	sm := tokenizer.NewTokenizerStateMachine()

	err := sm.ProcessData(reader)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sm.GetASTConstructor().HasOpenElements() {
		t.FailNow()
	}
	a := sm.GetASTConstructor().GetAST()

	variableConfig := make(map[string]interface{})
	variableConfig["someVar2"] = 123
	variableConfig["someVar1"] = map[string]interface{}{"hello": 123, "world": 223}
	variableConfig["key18"] = "key18+"
	variableConfig["value19"] = "somevalue"
	variableConfig["value20"] = []int{1, 2, 3}
	rs, err := interpreter.Interpret(a, variableConfig)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	fmt.Println(rs)

	var someMap map[string]interface{}
	err = json.Unmarshal([]byte(rs), &someMap)
	if err != nil {
		t.FailNow()
	}
	fmt.Printf("%#v", someMap)
	if len(someMap) != 11 {
		t.FailNow()
	}
}
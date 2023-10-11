package tokenizer_test

import (
	"bytes"
	"os"
	"runtime/pprof"
	"strings"
	"testing"

	"github.com/jaksonlin/go-jsonextend/tokenizer"

	_ "net/http/pprof"
)

func BenchmarkMyJsonE(b *testing.B) {
	cpuFile, err := os.Create("cpu.pprof")
	if err != nil {
		b.Fatal(err)
	}
	pprof.StartCPUProfile(cpuFile)

	defer pprof.StopCPUProfile()

	b.ReportAllocs()
	const sampleJson = `{
		"key1": "value1",
		"key2": 123,
		"key3": true,
		"key4": null,
		"key5": {
			"key6": "value6",
			"key7": 789,
			"key8": false,
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
		"key16": false,
		"key17": 123.4,
		"${key18}": "value18",
		"key19": "${value19}",
		"key20": ${value20}

	}`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(sampleJson)
		sm := tokenizer.NewTokenizerStateMachineFromIOReader(reader)

		b.StartTimer()
		err := sm.ProcessData()
		b.StopTimer()

		if err != nil {
			b.Error(err) // use b.Error to continue with other benchmarks
		}

		if sm.GetASTBuilder().HasOpenElements() {
			b.Error("Open element found")
		}
	}
	memFile, err := os.Create("mem.pprof")
	if err != nil {
		b.Fatal(err)
	}
	if err := pprof.WriteHeapProfile(memFile); err != nil {
		b.Fatal(err)
	}
	b.Log("end")
}

func TestMyJson1(t *testing.T) {
	const sampleJson = `[[1,2,3],[1,2,3]]`

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

}

func TestNull(t *testing.T) {
	data := []byte(`null`)

	sm := tokenizer.NewTokenizerStateMachineFromIOReader(bytes.NewReader(data))

	err := sm.ProcessData()
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sm.GetASTBuilder().HasOpenElements() {
		t.FailNow()
	}

}

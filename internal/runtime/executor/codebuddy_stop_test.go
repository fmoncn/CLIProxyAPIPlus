package executor

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestFixCodeBuddyPayloadStopStringBecomesArray(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":"</block>"}`)
	stop := gjson.GetBytes(fixCodeBuddyPayload(in), "stop")
	if !stop.IsArray() || len(stop.Array()) != 1 || stop.Array()[0].String() != "</block>" {
		t.Fatalf("stop = %s, want [\"</block>\"]", stop.Raw)
	}
}

func TestFixCodeBuddyPayloadStopArrayUnchanged(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":["a","b"]}`)
	stop := gjson.GetBytes(fixCodeBuddyPayload(in), "stop")
	if !stop.IsArray() || len(stop.Array()) != 2 {
		t.Fatalf("stop = %s, want [\"a\",\"b\"]", stop.Raw)
	}
}

func TestFixCodeBuddyPayloadNoStop(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}]}`)
	if gjson.GetBytes(fixCodeBuddyPayload(in), "stop").Exists() {
		t.Fatal("stop should not be added when absent")
	}
}

func TestFixCodeBuddyPayloadEmptyStopRemoved(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":""}`)
	if gjson.GetBytes(fixCodeBuddyPayload(in), "stop").Exists() {
		t.Fatal("empty string stop should be removed")
	}
}

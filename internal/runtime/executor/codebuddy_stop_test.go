package executor

import (
	"testing"

	"github.com/MoeForge/moego/ojson"
)

func TestFixCodeBuddyPayloadStopStringBecomesArray(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":"</block>"}`)
	stop := ojson.GetBytesNoCopy(fixCodeBuddyPayload(in), "stop")
	if !stop.IsArray() || len(stop.Array()) != 1 || stop.Array()[0].String() != "</block>" {
		t.Fatalf("stop = %s, want [\"</block>\"]", stop.Raw)
	}
}

func TestFixCodeBuddyPayloadStopArrayUnchanged(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":["a","b"]}`)
	stop := ojson.GetBytesNoCopy(fixCodeBuddyPayload(in), "stop")
	if !stop.IsArray() || len(stop.Array()) != 2 {
		t.Fatalf("stop = %s, want [\"a\",\"b\"]", stop.Raw)
	}
}

func TestFixCodeBuddyPayloadEmptyStopDropped(t *testing.T) {
	in := []byte(`{"messages":[{"role":"system","content":"s"}],"stop":""}`)
	stop := ojson.GetBytesNoCopy(fixCodeBuddyPayload(in), "stop")
	if stop.Exists() {
		t.Fatalf("stop = %s, want absent", stop.Raw)
	}
}

func TestFixCodeBuddyPayloadInjectsSystemWhenMissing(t *testing.T) {
	in := []byte(`{"messages":[{"role":"user","content":"hi"}]}`)
	out := fixCodeBuddyPayload(in)
	messages := ojson.GetBytesNoCopy(out, "messages")
	arr := messages.Array()
	if len(arr) != 2 {
		t.Fatalf("messages len = %d, want 2", len(arr))
	}
	if arr[0].Get("role").String() != "system" {
		t.Fatalf("first role = %s, want system", arr[0].Get("role").String())
	}
}

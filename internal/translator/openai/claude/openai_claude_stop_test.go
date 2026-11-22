package claude

import (
	"testing"

	"github.com/MoeForge/moego/ojson"
)

func TestConvertClaudeRequestToOpenAISingleStopIsArray(t *testing.T) {
	in := []byte(`{"model":"m","max_tokens":10,"stop_sequences":["</block>"],"messages":[{"role":"user","content":"hi"}]}`)
	out, err := ConvertClaudeRequestToOpenAI("m", in, false)
	if err != nil {
		t.Fatalf("ConvertClaudeRequestToOpenAI: %v", err)
	}
	stop := ojson.GetBytesNoCopy(out, "stop")
	if !stop.IsArray() || len(stop.Array()) != 1 || stop.Array()[0].String() != "</block>" {
		t.Fatalf("stop = %s, want [\"</block>\"]", stop.Raw)
	}
}

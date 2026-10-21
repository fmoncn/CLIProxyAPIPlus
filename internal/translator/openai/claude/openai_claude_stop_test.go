package claude

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestConvertClaudeRequestToOpenAISingleStopIsArray(t *testing.T) {
	in := []byte(`{"model":"m","max_tokens":10,"stop_sequences":["</block>"],"messages":[{"role":"user","content":"hi"}]}`)
	stop := gjson.GetBytes(ConvertClaudeRequestToOpenAI("m", in, false), "stop")
	if !stop.IsArray() || len(stop.Array()) != 1 || stop.Array()[0].String() != "</block>" {
		t.Fatalf("stop = %s, want [\"</block>\"]", stop.Raw)
	}
}

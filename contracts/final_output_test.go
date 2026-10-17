package contracts

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFinalOutputPolicyValidation(t *testing.T) {
	base := FinalOutputPolicy{Protocol: FinalOutputProtocol, Tool: "mcp_fixture_check", InputField: "draft", ReadyPath: []string{"structuredContent", "ready"}}
	if err := base.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"protocol", "tool", "reserved", "state_reserved", "field", "empty_path", "deep_path", "array_path", "url_path"} {
		t.Run(mode, func(t *testing.T) {
			p := base
			switch mode {
			case "protocol":
				p.Protocol = "unknown"
			case "tool":
				p.Tool = strings.Repeat("a", 129)
			case "reserved":
				p.Tool = FinalOutputToolName
			case "state_reserved":
				p.Tool = "working_state_update"
			case "field":
				p.InputField = "arguments.draft"
			case "empty_path":
				p.ReadyPath = nil
			case "deep_path":
				p.ReadyPath = make([]string, 9)
			case "array_path":
				p.ReadyPath = []string{"content", "0", "text"}
			case "url_path":
				p.ReadyPath = []string{"https://untrusted.invalid"}
			}
			if p.Validate() == nil {
				t.Fatal("accepted invalid policy")
			}
		})
	}
}

func TestFinalOutputAbsentPreservesHistoricalWire(t *testing.T) {
	var a Agent
	b, err := json.Marshal(a)
	if err != nil || strings.Contains(string(b), "final_output") {
		t.Fatal("changed historical JSON", err)
	}
	if err = json.Unmarshal([]byte(`{"final_output":{"protocol":"receipt-field-v1","tool":"check","input_field":"text","ready_path":["ready"]}}`), &a); err != nil || a.FinalOutput == nil || a.FinalOutput.Validate() != nil {
		t.Fatal(a, err)
	}
	if a.Validate() == nil {
		t.Fatal("accepted final output without native capabilities")
	}
}

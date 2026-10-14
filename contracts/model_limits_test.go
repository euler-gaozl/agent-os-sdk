package contracts

import (
	"encoding/json"
	"strings"
	"testing"
)

func nativeAgent() Agent {
	return Agent{Name: "fixture", Prompt: strings.Repeat("资料", 20000), ContextTokens: 1048576,
		Model: ModelConfig{BaseURL: "https://provider.example", Model: "fixture-model", MaxTokens: 393216,
			Limits: &ModelLimits{Model: "fixture-model", BaseURL: "https://provider.example", Protocol: "openai-chat-completions", ContextWindow: 1048576, MaxOutputTokens: 393216, Source: "https://provider.example/models",
				ObservedAt: "2026-10-14T00:00:00Z", MetadataSHA256: strings.Repeat("a", 64)}}}
}

func TestProviderCapabilitiesDoNotApplyLegacySmallAllowances(t *testing.T) {
	a := nativeAgent()
	if err := a.Validate(); err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Agent
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if err := decoded.Validate(); err != nil {
		t.Fatal(err)
	}
	if decoded.Model.MaxTokens != 393216 || decoded.ContextTokens != 1048576 || decoded.Model.Limits.MetadataSHA256 != a.Model.Limits.MetadataSHA256 {
		t.Fatal("provider capabilities were not preserved in snapshot JSON")
	}
	a.Model.Limits = nil
	if a.Validate() == nil {
		t.Fatal("legacy policy was silently widened")
	}
}

func TestProviderCapabilityObservationCannotFallBackOrDrift(t *testing.T) {
	for name, mutate := range map[string]func(*Agent){
		"unresolved":           func(a *Agent) { a.Model.Limits = &ModelLimits{} },
		"model drift":          func(a *Agent) { a.Model.Model = "different" },
		"endpoint drift":       func(a *Agent) { a.Model.BaseURL = "https://other.example" },
		"source drift":         func(a *Agent) { a.Model.Limits.Source = "https://other.example/models" },
		"wrong protocol":       func(a *Agent) { a.Model.Limits.Protocol = "other" },
		"context drift":        func(a *Agent) { a.ContextTokens = 64000 },
		"output drift":         func(a *Agent) { a.Model.MaxTokens = 6144 },
		"negative":             func(a *Agent) { a.Model.MaxTokens = -1; a.Model.Limits.MaxOutputTokens = -1 },
		"output beyond window": func(a *Agent) { a.Model.MaxTokens = 2000000; a.Model.Limits.MaxOutputTokens = 2000000 },
		"secret query":         func(a *Agent) { a.Model.Limits.Source += "?key=fixture" },
		"userinfo":             func(a *Agent) { a.Model.Limits.Source = "https://fixture@provider.example/models" },
		"missing time":         func(a *Agent) { a.Model.Limits.ObservedAt = "" },
		"missing digest":       func(a *Agent) { a.Model.Limits.MetadataSHA256 = "" },
	} {
		t.Run(name, func(t *testing.T) {
			a := nativeAgent()
			mutate(&a)
			if a.Validate() == nil {
				t.Fatal("invalid capability accepted")
			}
		})
	}
}

func TestLegacyModelJSONAndDefaultsRemainCompatible(t *testing.T) {
	a := Agent{Name: "fixture", Prompt: "fixture", Model: ModelConfig{BaseURL: "https://provider.example", Model: "fixture"}}
	if err := a.Validate(); err != nil {
		t.Fatal(err)
	}
	if a.ContextTokens != 32000 || a.Model.MaxTokens != 0 {
		t.Fatal("legacy defaults changed")
	}
	encoded, err := json.Marshal(a.Model)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "limits") {
		t.Fatal("new field entered legacy request JSON")
	}
}

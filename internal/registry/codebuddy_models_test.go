package registry

import "testing"

// The CodeBuddy catalog is split by account domain: domestic accounts are
// served models the international host does not have and vice versa. The
// built-ins below are the ones verified live against each host (2026-10-15
// probe); retired ids are deliberately absent because the upstream answers
// 11102 "service info not found" for them.
func TestGetCodeBuddyModelsIncludesVerifiedBuiltIns(t *testing.T) {
	models := GetCodeBuddyModels()
	byID := make(map[string]*ModelInfo, len(models))
	for _, model := range models {
		if model == nil {
			continue
		}
		if _, exists := byID[model.ID]; exists {
			t.Fatalf("duplicate CodeBuddy model ID %q", model.ID)
		}
		byID[model.ID] = model
	}

	for _, id := range []string{"hy3-preview", "kimi-k2.6", "deepseek-v4.1-flash", "glm-5.3"} {
		if byID[id] == nil {
			t.Fatalf("expected CodeBuddy model %q to be registered", id)
		}
	}

	if byID["hy3-preview"].Thinking == nil {
		t.Fatal("expected hy3-preview to advertise thinking support")
	}
}

func TestGetCodeBuddyModelsSplitByDomain(t *testing.T) {
	domestic := make(map[string]bool)
	for _, model := range GetCodeBuddyDomesticModels() {
		domestic[model.ID] = true
	}
	intl := make(map[string]bool)
	for _, model := range GetCodeBuddyIntlModels() {
		intl[model.ID] = true
	}
	for id := range domestic {
		if intl[id] {
			t.Fatalf("model %q registered on both domestic and international catalogs", id)
		}
	}
	if !intl["gpt-5.6-luna"] {
		t.Fatal("expected gpt-5.6-luna on the international catalog")
	}
	if !domestic["deepseek-v4.1-flash"] {
		t.Fatal("expected deepseek-v4.1-flash on the domestic catalog")
	}
}

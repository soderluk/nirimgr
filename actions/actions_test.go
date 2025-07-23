package actions

import (
	"testing"
)

func TestActionRegistryCreatesCorrectTypes(t *testing.T) {
	for name, model := range ActionRegistry {
		action := model()
		if action == nil {
			t.Errorf("ActionRegistry[%q] returned nil", name)
		}
		if action.GetName() != name {
			t.Errorf("ActionRegistry[%q] returned action with name %q", name, action.GetName())
		}
	}
}

func TestANameGetName(t *testing.T) {
	a := AName{Name: "TestAction"}
	if a.GetName() != "TestAction" {
		t.Errorf("AName.GetName() = %q, want %q", a.GetName(), "TestAction")
	}
}

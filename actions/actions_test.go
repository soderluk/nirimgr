package actions

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
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

type DummyAction struct {
	AName
	ID        uint64                `json:"id"`
	Reference WorkspaceReferenceArg `json:"reference"`
}

func TestSetActionId(t *testing.T) {
	ActionRegistry["dummy_action"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action"}} }
	defer delete(ActionRegistry, "dummy_action")

	// Prepare JSON for DummyAction
	payload := map[string]any{"id": 123, "reference": map[string]any{"id": 1, "index": 2, "name": "Foobar"}}
	d, err := json.Marshal(payload)
	assert.NoError(t, err)

	a := FromRegistry("dummy_action", d)
	a = SetActionID(a, 123)
	dummy, ok := a.(*DummyAction)
	assert.True(t, ok)
	assert.Equal(t, uint64(123), dummy.ID)
	assert.Equal(t, uint64(123), dummy.Reference.ID)
}

func TestFromRegistryMissingAction(t *testing.T) {
	// Should return nil if action type is not registered
	payload := map[string]any{"id": 1}
	d, err := json.Marshal(payload)
	assert.NoError(t, err)
	a := FromRegistry("nonexistent_action", d)
	assert.Nil(t, a)
}

func TestParseRawActions(t *testing.T) {
	// Register dummy action
	ActionRegistry["dummy_action"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action"}} }
	ActionRegistry["dummy_action_2"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action_2"}} }
	defer delete(ActionRegistry, "dummy_action")
	defer delete(ActionRegistry, "dummy_action_2")

	raw := map[string]json.RawMessage{
		"dummy_action":   []byte(`{"id": 42, "reference": {"id": 50}}`),
		"foo":            []byte(`{"id": 1}`),
		"dummy_action_2": []byte(`{"id": 99, "reference": {"id": 55}}`),
	}
	actions := ParseRawActions(raw)
	assert.Len(t, actions, 2)
	a1, ok1 := actions[0].(*DummyAction)
	a2, ok2 := actions[1].(*DummyAction)
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, uint64(42), a1.ID)
	assert.Equal(t, uint64(50), a1.Reference.ID)
	assert.Equal(t, uint64(99), a2.ID)
	assert.Equal(t, uint64(55), a2.Reference.ID)
}

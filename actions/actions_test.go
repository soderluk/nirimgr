package actions

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/soderluk/nirimgr/models"
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

type ExpectedDynamicID struct {
	ID, WindowID, ActiveWindowID, WorkspaceID, ReferenceID uint64
	ReferenceIndex, Index                                  uint8
	ReferenceName                                          string
}

func TestHandleDynamicIDs(t *testing.T) {
	ActionRegistry["dummy_action"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action"}} }
	defer delete(ActionRegistry, "dummy_action")

	// Prepare JSON for DummyAction
	payload := map[string]any{}
	d, err := json.Marshal(payload)
	assert.NoError(t, err)

	possibleKeyList := []struct {
		keys     models.PossibleKeys
		expected ExpectedDynamicID
	}{
		{
			keys:     models.PossibleKeys{ID: 123, WindowID: 456, ActiveWindowID: 789, WorkspaceID: 321, Index: 1, Reference: models.ReferenceKeys{ID: 123}},
			expected: ExpectedDynamicID{ID: 123, WindowID: 123, ActiveWindowID: 789, WorkspaceID: 321, Index: 1, ReferenceID: 123},
		},
		{
			keys:     models.PossibleKeys{WindowID: 456},
			expected: ExpectedDynamicID{WindowID: 456},
		},
		{
			keys:     models.PossibleKeys{ActiveWindowID: 789},
			expected: ExpectedDynamicID{ActiveWindowID: 789},
		},
		{
			keys:     models.PossibleKeys{WorkspaceID: 321},
			expected: ExpectedDynamicID{WorkspaceID: 321},
		},
		{
			keys:     models.PossibleKeys{Reference: models.ReferenceKeys{ID: 111}},
			expected: ExpectedDynamicID{ReferenceID: 111},
		},
		{
			keys:     models.PossibleKeys{Reference: models.ReferenceKeys{Index: 2}},
			expected: ExpectedDynamicID{ReferenceIndex: 2},
		},
		{
			keys:     models.PossibleKeys{Reference: models.ReferenceKeys{Name: "RefName"}},
			expected: ExpectedDynamicID{ReferenceName: "RefName"},
		},
	}
	for i, test := range possibleKeyList {
		a := FromRegistry("dummy_action", d)
		a = HandleDynamicIDs(a, test.keys)
		dummy, ok := a.(*DummyAction)
		assert.True(t, ok)
		assert.Equal(t, test.expected.ID, dummy.ID, fmt.Sprintf("%d: ID", i))
		assert.Equal(t, test.expected.ReferenceID, dummy.Reference.ID, fmt.Sprintf("%d: Reference.ID", i))
		assert.Equal(t, test.expected.ReferenceIndex, dummy.Reference.Index, fmt.Sprintf("%d: Reference.Index", i))
		assert.Equal(t, test.expected.ReferenceName, dummy.Reference.Name, fmt.Sprintf("%d: Reference.Name", i))
	}
}

func TestFromRegistryMissingAction(t *testing.T) {
	// Should return nil if action type is not registered
	payload := map[string]any{"id": 1}
	d, err := json.Marshal(payload)
	assert.NoError(t, err)
	a := FromRegistry("nonexistent_action", d)
	assert.Nil(t, a)
}

type ExpectedRawActions struct {
	ID, ReferenceID uint64
}

func TestParseRawActions(t *testing.T) {
	// Register dummy action
	ActionRegistry["dummy_action"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action"}} }
	ActionRegistry["dummy_action_2"] = func() Action { return &DummyAction{AName: AName{Name: "dummy_action_2"}} }
	defer delete(ActionRegistry, "dummy_action")
	defer delete(ActionRegistry, "dummy_action_2")

	tests := []struct {
		keys     map[string]json.RawMessage
		expected ExpectedRawActions
	}{
		{
			keys: map[string]json.RawMessage{
				"dummy_action": []byte(`{"id": 42, "reference": {"id": 50}}`),
			},
			expected: ExpectedRawActions{ID: 42, ReferenceID: 50},
		},
		{
			keys: map[string]json.RawMessage{
				"dummy_action_2": []byte(`{"id": 99, "reference": {"id": 55}}`),
			},
			expected: ExpectedRawActions{ID: 99, ReferenceID: 55},
		},
	}
	for i, test := range tests {
		actions := ParseRawActions(test.keys)
		assert.Len(t, actions, 1)
		dummy, ok := actions[0].(*DummyAction)
		assert.True(t, ok, "%d: not a DummyAction", i)
		assert.Equal(t, test.expected.ID, dummy.ID, "%d: ID", i)
		assert.Equal(t, test.expected.ReferenceID, dummy.Reference.ID, "%d: Reference.ID", i)
	}

	// Test with multiple actions in the map, order-independent
	raw := map[string]json.RawMessage{
		"dummy_action":   []byte(`{"id": 42, "reference": {"id": 50}}`),
		"foo":            []byte(`{"id": 1}`),
		"dummy_action_2": []byte(`{"id": 99, "reference": {"id": 55}}`),
	}
	actions := ParseRawActions(raw)
	assert.Len(t, actions, 2)
	found := map[uint64]*DummyAction{}
	for _, act := range actions {
		if d, ok := act.(*DummyAction); ok {
			found[d.ID] = d
		}
	}
	a1, ok1 := found[42]
	a2, ok2 := found[99]
	assert.True(t, ok1)
	assert.True(t, ok2)
	assert.Equal(t, uint64(50), a1.Reference.ID)
	assert.Equal(t, uint64(55), a2.Reference.ID)
}

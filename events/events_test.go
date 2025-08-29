package events

import (
	"encoding/json"
	"testing"

	"github.com/nalgeon/be"
	"github.com/soderluk/nirimgr/config"
	"github.com/soderluk/nirimgr/models"
)

type DummyEvent struct {
	EName
	Field string `json:"field"`
}

func TestParseEvent(t *testing.T) {
	EventRegistry["dummy_event"] = func() Event { return &DummyEvent{EName: EName{Name: "dummy_event"}} }
	defer delete(EventRegistry, "dummy_event")

	event := map[string]json.RawMessage{
		"dummy_event": json.RawMessage(`{"field": "value"}`),
	}
	key, model, err := ParseEvent(event)
	if err != nil {
		t.Fatalf("ParseEvent failed: %v", err)
	}
	if key != "dummy_event" {
		t.Errorf("Expected key 'dummy_event', got '%s'", key)
	}
	m, ok := model.(*DummyEvent)
	if !ok {
		t.Fatalf("Expected *DummyEvent, got %T", model)
	}
	if m.Field != "value" {
		t.Errorf("Expected field 'value', got '%s'", m.Field)
	}
}

func TestParseEvent_NoEventFound(t *testing.T) {
	event := map[string]json.RawMessage{
		"unknown_event": json.RawMessage(`{"field": "value"}`),
	}
	key, model, err := ParseEvent(event)
	if err == nil || err.Error() != "no event found" {
		t.Errorf("Expected error 'no event found', got %v", err)
	}
	if key != "" || model != nil {
		t.Errorf("Expected key '', model nil, got key '%s', model %v", key, model)
	}
}

func TestParseEvent_UnmarshalError(t *testing.T) {
	EventRegistry["dummy_event"] = func() Event { return &DummyEvent{EName: EName{Name: "dummy_event"}} }
	defer delete(EventRegistry, "dummy_event")

	event := map[string]json.RawMessage{
		"dummy_event": []byte(`{"field": }`),
	}
	key, model, err := ParseEvent(event)
	if err == nil {
		t.Errorf("Expected unmarshal error, got nil")
	}
	if key != "" {
		t.Errorf("Expected key 'dummy_event', got '%s'", key)
	}
	if model != nil {
		t.Errorf("Expected model nil, got %v", model)
	}
}

func TestUpdateWindowMatched_MatchAndAction(t *testing.T) {
	window := &models.Window{ID: 1, Title: "Test window", AppID: "test-app"}
	existingWindows := map[uint64]*models.Window{}

	// Simulate config
	cfg := &models.Config{
		Rules: []models.Rule{
			{
				Match: []models.Match{
					{
						AppID: "test-app",
					},
				},
			},
		},
	}
	config.Config = cfg
	matchWindowAndPerformActions(window, existingWindows)
	if !window.Matched {
		t.Errorf("Expected window to be matched")
	}
}

func TestUpdateWorkspaceMatched_MatchAndAction(t *testing.T) {
	workspace := &models.Workspace{ID: 1, Name: "Test workspace", Output: "test-output"}
	existingWorkspaces := map[uint64]*models.Workspace{}

	// Simulate config
	cfg := &models.Config{
		Rules: []models.Rule{
			{
				Type: "workspace",
				Match: []models.Match{
					{
						Name: "Test workspace",
					},
				},
			},
		},
	}
	config.Config = cfg
	matchWorkspaceAndPerformActions(workspace, existingWorkspaces)
	if !workspace.Matched {
		t.Errorf("Expected workspace to be matched")
	}
}

func TestEvaluateCondition(t *testing.T) {
	t.Run("pass", func(t *testing.T) {
		event := WindowClosed{EName: EName{"WindowClosed"}, ID: 1}
		got, _ := EvaluateCondition("event.ID == 1", event)
		be.True(t, got)
		got, _ = EvaluateCondition("event.ID == 2", event)
		be.True(t, !got)
	})
	t.Run("with error", func(t *testing.T) {
		event := WindowFocusChanged{EName: EName{"WindowFocusChanged"}, ID: 2}
		_, err := EvaluateCondition("event.Urgent == true", event)
		be.Err(t, err)
	})
}

package actions

import (
	"encoding/json"
	"testing"
)

func TestColumnDisplay(t *testing.T) {
	cd := ColumnDisplay{Normal: 1, Tabbed: 2}
	b, err := json.Marshal(cd)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want := `{"Normal":1,"Tabbed":2}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}

	// Test omitempty
	cd = ColumnDisplay{Normal: 0, Tabbed: 0}
	b, err = json.Marshal(cd)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want = `{}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}
}

func TestLayoutSwitchTarget(t *testing.T) {
	lst := LayoutSwitchTarget{Next: 1, Prev: 2, Index: 3}
	b, err := json.Marshal(lst)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want := `{"Next":1,"Prev":2,"Index":3}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}

	// Test omitempty
	lst = LayoutSwitchTarget{}
	b, err = json.Marshal(lst)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want = `{}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}
}

func TestPositionChange(t *testing.T) {
	pc := PositionChange{SetFixed: 10.5, AdjustFixed: -2.5}
	b, err := json.Marshal(pc)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want := `{"SetFixed":10.5,"AdjustFixed":-2.5}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}

	// Test omitempty
	pc = PositionChange{}
	b, err = json.Marshal(pc)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want = `{}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}
}

func TestSizeChange(t *testing.T) {
	sc := SizeChange{SetFixed: 100, SetProportion: 0.5, AdjustFixed: -20, AdjustProportion: 0.1}
	b, err := json.Marshal(sc)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want := `{"SetFixed":100,"SetProportion":0.5,"AdjustFixed":-20,"AdjustProportion":0.1}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}

	// Test omitempty
	sc = SizeChange{}
	b, err = json.Marshal(sc)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want = `{}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}
}

func TestWorkspaceReferenceArg(t *testing.T) {
	wra := WorkspaceReferenceArg{ID: 42, Index: 2, Name: "main"}
	b, err := json.Marshal(wra)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want := `{"Id":42,"Index":2,"Name":"main"}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}

	// Test omitempty
	wra = WorkspaceReferenceArg{}
	b, err = json.Marshal(wra)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	want = `{}`
	if string(b) != want {
		t.Errorf("expected %s, got %s", want, string(b))
	}
}

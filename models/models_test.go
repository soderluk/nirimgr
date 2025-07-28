package models

import (
	"testing"
)

var testConfig = &Config{
	Rules: []Rule{
		{
			Type: "window",
			Match: []Match{
				{Title: "Bitwarden", AppID: "zen"},
			},
			Exclude: nil,
		},
		{
			Type: "window",
			Match: []Match{
				{AppID: "^foot$"},
				{AppID: "^mpv$"},
			},
		},
		{
			Type: "window",
			Match: []Match{
				{Title: "^foo$"},
			},
			Exclude: []Match{
				{AppID: "^bar$"},
			},
		},
		{
			Type: "window",
			Match: []Match{
				{Title: ""},
			},
		},
		{
			Type: "window",
			Match: []Match{
				{AppID: ""},
			},
		},
		{
			Type: "window",
			Match: []Match{
				{Title: "", AppID: ""},
			},
		},
		{
			Type: "workspace",
			Match: []Match{
				{Name: "chat"},
			},
		},
		{
			Type: "workspace",
			Match: []Match{
				{Name: "work", Output: "eDP-1"},
			},
		},
		{
			Type: "workspace",
			Match: []Match{
				{Name: ""},
			},
		},
		{
			Type: "workspace",
			Match: []Match{
				{Output: ""},
			},
		},
		{
			Type: "workspace",
			Match: []Match{
				{Output: "", Name: ""},
			},
		},
	},
}

func TestWindowRules(t *testing.T) {
	windows := []Window{
		{
			ID:    1,
			Title: "Bitwarden",
			AppID: "zen",
		},
		{
			ID:    2,
			AppID: "foot",
		},
		{
			ID:    3,
			Title: "mpv",
		},
		{
			ID:    4,
			Title: "foo",
			AppID: "bar",
		},
	}

	tests := []struct {
		ruleIdx   int
		windowIdx int
		wantMatch bool
	}{
		{0, 0, true},
		{0, 1, false},
		{1, 1, true},
		{1, 2, false},
		{2, 3, false},
		{3, 0, false},
		{4, 0, false},
		{5, 0, false},
	}
	rules := testConfig.GetRules()
	for _, tt := range tests {
		rule := rules[tt.ruleIdx]
		window := windows[tt.windowIdx]
		got := rule.WindowMatches(window)
		if got != tt.wantMatch {
			t.Errorf("Rule[%d].Matches(Window[%d]) = %v, want %v", tt.ruleIdx, tt.windowIdx, got, tt.wantMatch)
		}
	}
}

func TestWorkspaceRules(t *testing.T) {
	workspaces := []Workspace{
		{
			ID:     1,
			Idx:    1,
			Name:   "chat",
			Output: "DP-8",
		},
		{
			ID:     2,
			Idx:    2,
			Name:   "work",
			Output: "eDP-1",
		},
		{
			ID:  3,
			Idx: 3,
		},
	}

	tests := []struct {
		ruleIdx      int
		workspaceIdx int
		wantMatch    bool
	}{
		{6, 0, true},
		{6, 1, false},
		{6, 2, false},
		{7, 0, false},
		{7, 1, true},
		{7, 2, false},
		{8, 2, false},
		{9, 2, false},
		{10, 2, false},
	}

	rules := testConfig.GetRules()
	for _, tt := range tests {
		rule := rules[tt.ruleIdx]
		workspace := workspaces[tt.workspaceIdx]
		got := rule.WorkspaceMatches(workspace)
		if got != tt.wantMatch {
			t.Errorf("Rule[%d].Matches(Workspace[%d]) = %v, want %v", tt.ruleIdx, tt.workspaceIdx, got, tt.wantMatch)
		}
	}
}

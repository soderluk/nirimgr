package models

import (
	"testing"
)

func TestRuleMatches(t *testing.T) {
	rules := []Rule{
		{
			Match: []Match{
				{Title: "Bitwarden", AppID: "zen"},
			},
			Exclude: nil,
		},
		{
			Match: []Match{
				{AppID: "^foot$"},
				{AppID: "^mpv$"},
			},
		},
		{
			Match: []Match{
				{Title: "^foo$"},
			},
			Exclude: []Match{
				{AppID: "^bar$"},
			},
		},
	}
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
	}

	for _, tt := range tests {
		rule := rules[tt.ruleIdx]
		window := windows[tt.windowIdx]
		got := rule.Matches(window)
		if got != tt.wantMatch {
			t.Errorf("Rule[%d].Matches(Window[%d]) = %v, want %v", tt.ruleIdx, tt.windowIdx, got, tt.wantMatch)
		}
	}
}

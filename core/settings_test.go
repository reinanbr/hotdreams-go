package core

import (
	"fmt"
	"testing"
)

func TestPyQuote(t *testing.T) {
	cases := map[string]string{
		"lorena aquino": "lorena%20aquino",
		"milf big ass":  "milf%20big%20ass",
		"a/b":           "a/b",
		"":              "",
	}
	for in, want := range cases {
		if got := pyQuote(in); got != want {
			t.Errorf("pyQuote(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSearchPornBaseStopsOnEmptyPage(t *testing.T) {
	calls := 0
	fetch := func(url string, page int, query string) ([]VideoData, bool, error) {
		calls++
		if page <= 2 {
			return []VideoData{{Title: fmt.Sprintf("v%d", page)}}, true, nil
		}
		return nil, false, nil
	}
	build := func(query string, page int) string {
		return fmt.Sprintf("http://x/%d", page)
	}

	data, err := SearchPornBase("q", "http://x", "site", 5, nil, fetch, build)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 fetch calls (2 hits + 1 stop), got %d", calls)
	}
	if data.LenVideos != 2 {
		t.Errorf("expected 2 videos, got %d", data.LenVideos)
	}
	if data.LenPages == nil || *data.LenPages != 2 {
		t.Errorf("expected LenPages=2, got %v", data.LenPages)
	}
	if !data.Success {
		t.Error("expected Success=true")
	}
}

func TestSearchPornBaseSinglePage(t *testing.T) {
	fetch := func(url string, page int, query string) ([]VideoData, bool, error) {
		return []VideoData{{Title: "only"}}, true, nil
	}
	build := func(query string, page int) string { return "http://x" }

	page := 3
	data, err := SearchPornBase("q", "http://x", "site", 0, &page, fetch, build)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.LenPages != nil {
		t.Errorf("expected LenPages=nil for single-page search, got %v", *data.LenPages)
	}
	if data.LenVideos != 1 {
		t.Errorf("expected 1 video, got %d", data.LenVideos)
	}
}

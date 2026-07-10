package xvideos

import (
	"strings"
	"testing"

	"github.com/reinanbr/hotdreams-go/core"
)

const searchPageHTML = `
<html><body>
<div id="video_12345">
  <a href="/video12345/some_title"></a>
  <a href="/video12345/some_title">Delicious busty loves blowjob</a>
  <span class="duration">11 min</span>
  <img data-src="https://img-cf.xvideos-cdn.com/thumb.jpg">
</div>
</body></html>`

func TestMinerateVideo(t *testing.T) {
	doc, err := core.ParseHTML([]byte(searchPageHTML))
	if err != nil {
		t.Fatalf("ParseHTML: %v", err)
	}
	videos, ok := minerateVideo(doc, 1, "https://www.xvideos.com/?k=x&p=1")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(videos))
	}
	v := videos[0]
	if v.Title != "Delicious busty loves blowjob" {
		t.Errorf("Title = %q", v.Title)
	}
	if v.URL != "https://www.xvideos.com/video12345/some_title" {
		t.Errorf("URL = %q", v.URL)
	}
	if v.Time != "11 min" {
		t.Errorf("Time = %q", v.Time)
	}
}

func TestMinerateVideoNoResults(t *testing.T) {
	doc, _ := core.ParseHTML([]byte(`<html><body>No video match with this search</body></html>`))
	_, ok := minerateVideo(doc, 1, "url")
	if ok {
		t.Fatal("expected ok=false on no-results sentinel")
	}
}

func TestEmbedKeyExtraction(t *testing.T) {
	url := "https://www.xvideos.com/video58733211/gali_diva_hot_sex_with_her_boyfriend"
	segments := strings.Split(url, "video")
	if len(segments) < 3 {
		t.Fatalf("unexpected split result: %#v", segments)
	}
	key := strings.Split(segments[2], "/")[0]
	if key != "58733211" {
		t.Errorf("key = %q, want 58733211", key)
	}
}

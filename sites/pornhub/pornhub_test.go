package pornhub

import (
	"testing"

	"github.com/reinanbr/hot_dreams/core"
)

const searchPageHTML = `
<html><body>
<ul>
<li class="pcVideoListItem">
  <a href="/view_video.php?viewkey=abc123">
    <img alt="Some Video Title" src="https://di.phncdn.com/thumb.jpg" data-mediabook="https://dw.phncdn.com/preview.webm">
  </a>
  <var>12:37</var>
  <var class="added">6 days ago</var>
  <span class="views">14.7M views</span>
  <div class="rating-container">89 %</div>
</li>
</ul>
</body></html>`

func TestMinerateVideo(t *testing.T) {
	doc, err := core.ParseHTML([]byte(searchPageHTML))
	if err != nil {
		t.Fatalf("ParseHTML: %v", err)
	}
	videos, ok := minerateVideo(doc, 1, "https://pornhub.com/video/search?search=x&page=1")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(videos))
	}
	v := videos[0]
	if v.Title != "Some Video Title" {
		t.Errorf("Title = %q", v.Title)
	}
	if v.URL != "https://pornhub.com/view_video.php?viewkey=abc123" {
		t.Errorf("URL = %q", v.URL)
	}
	if v.Time != "12:37" {
		t.Errorf("Time = %q", v.Time)
	}
	if v.DateUpload != "6 days ago" {
		t.Errorf("DateUpload = %q", v.DateUpload)
	}
	if v.Rating != "89%" {
		t.Errorf("Rating = %q", v.Rating)
	}
	if v.SiteName != SiteName {
		t.Errorf("SiteName = %q", v.SiteName)
	}
}

func TestMinerateVideoEndOfSearch(t *testing.T) {
	doc, _ := core.ParseHTML([]byte(`<html><body>We're sorry, but the requested search cannot be found. Broaden your search.</body></html>`))
	_, ok := minerateVideo(doc, 1, "url")
	if ok {
		t.Fatal("expected ok=false on end-of-search sentinel")
	}
}

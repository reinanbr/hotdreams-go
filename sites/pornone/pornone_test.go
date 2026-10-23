package pornone

import (
	"testing"

	"github.com/reinanbr/hotdreams-go/core"
)

func TestVideoCardsLazyThumbnail(t *testing.T) {
	doc, err := core.ParseHTML([]byte(`<div class="mt-1">
<a class="videocard" href="/blonde/some-video/123456/">
  <img class="imgvideo" src="data:image/gif;base64,R0lGODlhAQABAAAAACw=" data-src="https://th-eu3.pornone.com/t/12/123456/b1.jpg">
  <span class="durlabel">12:34</span>
  <div class="videotitle">Some video</div>
  <div class="author">user  1.2K views</div>
</a></div>`))
	if err != nil {
		t.Fatal(err)
	}
	videos := videoCards(doc.Find("div.mt-1"), 1, URLBase)
	if len(videos) != 1 {
		t.Fatalf("expected 1 video, got %d", len(videos))
	}
	v := videos[0]
	if v.Thumbnail != "https://th-eu3.pornone.com/t/12/123456/b1.jpg" {
		t.Errorf("thumbnail = %q", v.Thumbnail)
	}
	if v.URL != "https://pornone.com/blonde/some-video/123456/" {
		t.Errorf("url = %q", v.URL)
	}
	if v.Duration != 754 {
		t.Errorf("duration = %d", v.Duration)
	}
}

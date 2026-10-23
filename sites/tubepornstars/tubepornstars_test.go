package tubepornstars

import (
	"testing"

	"github.com/reinanbr/hotdreams-go/core"
)

const outLink = "/out/?l=3AASHM4hehSbq2txY29aN0ZNZmRzAtmbaHR0cHM6Ly93d3cucHhvaHViLmNvbS92aWRlb3MvMjk2NzUyNy9uYXR1cmFsLWJvb2JzLWNodWJieS1taWxmLW5hdGFzaGEtaGVscHMtaGVyLWxvc3Qtc3RlcHNvbi1uYXRhc2hhLW5pY2UvP3V0bV9zb3VyY2U9YXduJnV0bV9tZWRpdW09dGdwJnV0bV9jYW1wYWlnbj1jcGPNAZCidGMBNKdwb3B1bGFyBNkweyJhbGwiOiIiLCJvcmllbnRhdGlvbiI6InN0cmFpZ2h0IiwicHJpY2luZyI6IiJ9zXLyzmrbf%2BOocG9ybnN0YXLNxi3A2XxbeyIxIjoielRaRGdpcVFXd2oifSx7IjIiOiJ3ZXJaZ2QxdTM1VyJ9LHsiMyI6ImZ4aWlxaVdFZkhhIn0seyItMSI6ImdRQWlpUkZ6amI1In0seyItMiI6IkVFTFNJdDM0QnVBIn0seyItMyI6IktEbnMxZnBCdFlTIn1d&amp;c=83fcf75d91f953ff&amp;v=3&amp;"

const wantURL = "https://www.pxohub.com/videos/2967527/natural-boobs-chubby-milf-natasha-helps-her-lost-stepson-natasha-nice/?utm_source=awn&utm_medium=tgp&utm_campaign=cpc"

const searchPageHTML = `<html><body>
<div class="card sub group relative block space-y-1 paid" data-public-id="ad">
  <a class="item-title item-link" href="https://ads.example/x" title="Ad">Ad</a>
</div>
<div class="card sub group relative block space-y-1" data-public-id="kqcoZ7FMfds">
  <a class="item-link rate-link" href="` + outLink + `">
    <img class="item-image size-full" src="https://c4.ttcache.com/thumbnail/kqcoZ7FMfds/288x162/a.jpg"/>
    <span class="item-meta-container">
      <span class="rating-badge"><span class="item-score text-sm"><i class="far fa-thumbs-up"></i>63%</span></span>
      <span class="badge float-end"><span class="font-bold italic">
          HD
        </span>
        31:06 </span>
    </span>
  </a>
  <div class="item-footer"><div class="relative w-full mt-1">
    <a class="item-title item-link rate-link" href="` + outLink + `" title="Natural Boobs Chubby MILF Natasha Helps Her Lost Stepson - Natasha Nice"><span>Natural Boobs</span></a>
    <div class="item-source-rating-container mt-1">
      <a class="item-source block" href="/source/pornoxo"><i class="far fa-badge-check"></i>PornoXO
      </a><span class="item-source block">4 days ago</span>
    </div>
  </div></div>
</div>
</body></html>`

func TestSourceURL(t *testing.T) {
	if got := sourceURL(core.ResolveURL(URLBase, "/out/?l=3AASHM4hehSbq2txY29aN0ZNZmRzAtmbaHR0cHM6Ly93d3cucHhvaHViLmNvbS92aWRlb3MvMjk2NzUyNy9uYXR1cmFsLWJvb2JzLWNodWJieS1taWxmLW5hdGFzaGEtaGVscHMtaGVyLWxvc3Qtc3RlcHNvbi1uYXRhc2hhLW5pY2UvP3V0bV9zb3VyY2U9YXduJnV0bV9tZWRpdW09dGdwJnV0bV9jYW1wYWlnbj1jcGPNAZCidGM%3D&c=1")); got != wantURL {
		t.Errorf("sourceURL = %q", got)
	}
	if got := sourceURL("/out/?l=not-base64!"); got != "" {
		t.Errorf("sourceURL(garbage) = %q, want empty", got)
	}
}

func TestMinerateVideo(t *testing.T) {
	doc, err := core.ParseHTML([]byte(searchPageHTML))
	if err != nil {
		t.Fatal(err)
	}
	videos, ok := minerateVideo(doc, 1, "url")
	if !ok {
		t.Fatal("expected ok=true")
	}
	if len(videos) != 1 {
		t.Fatalf("expected 1 video (ad skipped), got %d", len(videos))
	}
	v := videos[0]
	if v.Title != "Natural Boobs Chubby MILF Natasha Helps Her Lost Stepson - Natasha Nice" {
		t.Errorf("Title = %q", v.Title)
	}
	if v.URL != wantURL {
		t.Errorf("URL = %q", v.URL)
	}
	if v.Time != "31:06" || v.Duration != 1866 {
		t.Errorf("Time = %q, Duration = %d", v.Time, v.Duration)
	}
	if v.Thumbnail != "https://c4.ttcache.com/thumbnail/kqcoZ7FMfds/288x162/a.jpg" {
		t.Errorf("Thumbnail = %q", v.Thumbnail)
	}
	if v.Rating != "63%" || v.RatingInt != 63 {
		t.Errorf("Rating = %q, RatingInt = %d", v.Rating, v.RatingInt)
	}
	if v.Stats != "PornoXO" || v.DateUpload != "4 days ago" {
		t.Errorf("Stats = %q, DateUpload = %q", v.Stats, v.DateUpload)
	}
}

func TestMinerateVideoNoResults(t *testing.T) {
	doc, _ := core.ParseHTML([]byte(`<html><body><h3>No results found for the following search: <span>Zzq</span></h3>
<div class="card sub"><a class="item-title" href="/out/?l=x" title="fallback"></a></div></body></html>`))
	if _, ok := minerateVideo(doc, 1, "url"); ok {
		t.Fatal("expected ok=false on no-results page")
	}
}

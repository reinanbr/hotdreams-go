// Package playvids ports hotdreams.sites.playvids to Go.
//
// Unlike every other site in this lib, playvids.search_porn in Python returns
// a plain list of dicts (not a DataVideos of VideoData) and never went through
// search_porn_base — this port keeps that same shape (a []Video, no pagination
// metadata) for fidelity.
package playvids

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "playvids"
	URLBase    = "https://www.playvids.com"
	cookieFile = ".cookies_playvids.json"
)

var client = core.NewClient(cookieFile, URLBase)

// Video mirrors the plain dict hotdreams.sites.playvids.search_porn returns
// in Python: {'title','duration','imgUrl','url','dur'}.
type Video struct {
	Title    string
	Duration string
	ImgURL   string
	URL      string
	Dur      int
}

// SearchPorn searches playvids for query, fetching up to pageLimit pages
// (default 2). If long is true, restricts to long-duration videos sorted by
// relevance, same as the Python `long=True` flag.
func SearchPorn(query string, pageLimit int, long bool) ([]Video, error) {
	if pageLimit <= 0 {
		pageLimit = 2
	}
	q := strings.Join(strings.Split(query, " "), "+")

	var videos []Video
	p := 1
	for i := 0; i < pageLimit; i++ {
		var url string
		if long {
			url = fmt.Sprintf("%s/videos?q=%s&duration=long&sort_by=relevance&page=%d", URLBase, q, p)
		} else {
			url = fmt.Sprintf("%s/videos?q=%s&page=%d", URLBase, q, p)
		}

		body, _, _, err := client.Fetch(url)
		if err != nil {
			return videos, err
		}
		doc, err := core.ParseHTML(body)
		if err != nil {
			return videos, err
		}
		text := doc.Text()
		if strings.Contains(text, "did not match any videos.") || strings.Contains(text, "Page not found") {
			break
		}
		p++

		doc.Find("div").Each(func(_ int, div *goquery.Selection) {
			if !div.HasClass("itemVideo") {
				return
			}
			durLabel := div.Find("span.duration").First().Text()
			var durSeconds int
			parts := strings.SplitN(durLabel, ":", 2)
			if len(parts) == 2 {
				mins, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
				secs, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
				if err1 == nil && err2 == nil {
					durSeconds = secs + 60*mins
				}
			}
			title, _ := div.Find("img").First().Attr("alt")
			imgURL, _ := div.Find("img").First().Attr("src")
			href, _ := div.Find("a").First().Attr("href")

			videos = append(videos, Video{
				Title:    title,
				Duration: durLabel,
				ImgURL:   imgURL,
				URL:      URLBase + href,
				Dur:      durSeconds,
			})
		})
	}
	return videos, nil
}

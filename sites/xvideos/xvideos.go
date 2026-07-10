// Package xvideos ports hotdreams.sites.xvideos to Go.
package xvideos

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "XVideos"
	URLBase    = "https://www.xvideos.com"
	cookieFile = ".cookies_XVideos.json"
)

var client = core.NewClient(cookieFile, URLBase)

// minerateVideo mirrors xvideos/utils.py's minerate_video.
func minerateVideo(doc *goquery.Document, pageNumber int, urlFont string) ([]core.VideoData, bool) {
	if strings.Contains(doc.Text(), "No video match with this search") {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div[id]").Each(func(_ int, div *goquery.Selection) {
		id, _ := div.Attr("id")
		if !strings.Contains(id, "video") {
			return
		}

		links := div.Find("a")
		var title string
		if links.Length() > 1 {
			title = strings.TrimSpace(links.Eq(1).Text())
		} else {
			title = strings.TrimSpace(links.First().Text())
		}
		timeVal := div.Find("span.duration").First().Text()
		thumbnail, _ := div.Find("img").First().Attr("data-src")
		href, _ := links.First().Attr("href")

		videos = append(videos, core.VideoData{
			URL:        URLBase + href,
			Title:      title,
			Time:       timeVal,
			Thumbnail:  thumbnail,
			PageNumber: pageNumber,
			SiteName:   SiteName,
			URLFont:    urlFont,
			Indice:     i,
		})
		i++
	})
	return videos, true
}

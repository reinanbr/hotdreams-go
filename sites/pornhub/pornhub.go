// Package pornhub ports hotdreams.sites.pornhub to Go.
package pornhub

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

const (
	SiteName   = "PornHub"
	URLBase    = "https://pornhub.com"
	cookieFile = ".cookies_PornHub.json"
)

var client = core.NewClient(cookieFile, URLBase)

// minerateVideo mirrors pornhub/utils.py's minerate_video.
func minerateVideo(doc *goquery.Document, pageNumber int, urlFont string) ([]core.VideoData, bool) {
	if strings.Contains(doc.Text(), "We're sorry, but the requested search cannot be found. Broaden your search.") {
		return nil, false
	}

	var videos []core.VideoData
	doc.Find("li.pcVideoListItem").Each(func(i int, li *goquery.Selection) {
		img := li.Find("img").First()
		title, _ := img.Attr("alt")
		timeVal := li.Find("var").First().Text()
		thumbnail, _ := img.Attr("src")
		href, _ := li.Find("a").First().Attr("href")
		preview, _ := img.Attr("data-mediabook")
		views := li.Find("span.views").First().Text()
		rating := strings.Join(strings.Fields(li.Find("div.rating-container").First().Text()), "")
		uploadDate := li.Find("var.added").First().Text()

		videos = append(videos, core.VideoData{
			URL:        URLBase + href,
			Title:      title,
			Time:       timeVal,
			Thumbnail:  thumbnail,
			Preview:    preview,
			Views:      views,
			Rating:     rating,
			DateUpload: uploadDate,
			PageNumber: pageNumber,
			SiteName:   SiteName,
			URLFont:    urlFont,
			Indice:     i,
		})
	})
	return videos, true
}

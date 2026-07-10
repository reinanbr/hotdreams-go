// Package fpo ports hotdreams.sites.fpo (fpo.xxx) to Go.
package fpo

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "FPO.XXX"
	URLBase    = "https://www.fpo.xxx"
	cookieFile = ".cookies_fpo.json"
)

var client = core.NewClient(cookieFile, URLBase)

func isVideoCard(div *goquery.Selection) bool {
	linkWithHref := div.Find("a[href]").First()
	title := div.Find("strong.title").First()
	return linkWithHref.Length() > 0 && title.Length() > 0
}

// minerateVideo mirrors fpo/utils.py's minerate_video. urlFont is the
// request URL passed in by the caller (not the post-redirect URL), same as
// the Python code.
func minerateVideo(doc *goquery.Document, urlFont string, pageNumber int) ([]core.VideoData, bool) {
	var cards []*goquery.Selection
	doc.Find("div.item").Each(func(_ int, d *goquery.Selection) {
		if isVideoCard(d) {
			cards = append(cards, d)
		}
	})
	if len(cards) == 0 {
		return nil, false
	}

	var videos []core.VideoData
	for i, video := range cards {
		linkA := video.Find("a[href]").First()
		urlVideo, _ := linkA.Attr("href")
		imgTag := linkA.Find("img").First()

		titleTag := linkA.Find("strong.title").First()
		var title string
		if titleTag.Length() > 0 {
			title = strings.TrimSpace(titleTag.Text())
		} else {
			title, _ = imgTag.Attr("alt")
		}
		thumbnail, ok := imgTag.Attr("data-original")
		if !ok {
			thumbnail, _ = imgTag.Attr("src")
		}
		preview, _ := imgTag.Attr("data-preview")

		timeVal := strings.TrimSpace(linkA.Find("span.duration").First().Text())

		rating := strings.TrimSpace(linkA.Find("div.rating").First().Text())
		var ratingInt int
		if rating != "" {
			ratingInt, _ = strconv.Atoi(strings.ReplaceAll(rating, "%", ""))
		}

		views := strings.TrimSpace(linkA.Find("div.views").First().Text())
		var viewsInt int64
		if views != "" {
			cleaned := strings.ReplaceAll(strings.ReplaceAll(views, " ", ""), ",", "")
			if n, err := strconv.ParseInt(cleaned, 10, 64); err == nil {
				viewsInt = n
			}
		}

		dur, _ := core.ParseDurationClock(timeVal)

		videos = append(videos, core.VideoData{
			Title:           title,
			Time:            timeVal,
			DurationSeconds: dur,
			PageNumber:      pageNumber,
			URL:             urlVideo,
			URLFont:         urlFont,
			Views:           views,
			ViewsInt:        viewsInt,
			Rating:          rating,
			RatingInt:       ratingInt,
			Thumbnail:       thumbnail,
			SiteName:        SiteName,
			Preview:         preview,
			Indice:          i,
		})
	}
	return videos, true
}

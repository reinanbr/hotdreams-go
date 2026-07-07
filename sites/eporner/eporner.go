// Package eporner ports hotdreams.sites.eporner to Go.
package eporner

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

const (
	SiteName   = "Eporner"
	URLBase    = "https://www.eporner.com"
	cookieFile = ".cookies_eporner.json"
)

var client = core.NewClient(cookieFile, URLBase)

// minerateVideo mirrors eporner/utils.py's minerate_video. urlFont is the
// request URL passed in by the caller (not the post-redirect URL), same as
// the Python code.
func minerateVideo(doc *goquery.Document, urlFont string, pageNumber int) ([]core.VideoData, bool) {
	videoDiv := doc.Find("div.mb")
	if videoDiv.Length() == 0 || strings.Contains(doc.Text(), "No results") {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	videoDiv.Each(func(_ int, video *goquery.Selection) {
		thumbA := video.Find("a").First()
		href, hasHref := thumbA.Attr("href")
		if thumbA.Length() == 0 || !hasHref {
			return
		}
		imgTag := thumbA.Find("img").First()

		titleTag := video.Find("p.mbtit").First()
		var title string
		if titleTag.Length() > 0 {
			title = strings.TrimSpace(titleTag.Text())
		} else {
			title, _ = imgTag.Attr("alt")
		}
		thumbnail, ok := imgTag.Attr("data-src")
		if !ok {
			thumbnail, _ = imgTag.Attr("src")
		}

		timeVal := strings.TrimSpace(video.Find("span.mbtim").First().Text())

		rating := strings.TrimSpace(video.Find("span.mbrate").First().Text())
		var ratingInt int
		if rating != "" {
			ratingInt, _ = strconv.Atoi(strings.ReplaceAll(rating, "%", ""))
		}

		views := strings.TrimSpace(video.Find("span.mbvie").First().Text())
		var viewsInt int64
		if views != "" {
			if n, err := strconv.ParseInt(strings.ReplaceAll(views, ",", ""), 10, 64); err == nil {
				viewsInt = n
			}
		}

		dur, _ := core.ParseDurationClock(timeVal)

		videos = append(videos, core.VideoData{
			Title:           title,
			Time:            timeVal,
			DurationSeconds: dur,
			PageNumber:      pageNumber,
			URL:             URLBase + href,
			URLFont:         urlFont,
			Views:           views,
			ViewsInt:        viewsInt,
			Rating:          rating,
			RatingInt:       ratingInt,
			Thumbnail:       thumbnail,
			SiteName:        SiteName,
			Indice:          i,
		})
		i++
	})
	return videos, true
}

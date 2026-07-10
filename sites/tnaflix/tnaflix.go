// Package tnaflix ports hotdreams.sites.tnaflix to Go.
package tnaflix

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "tnaflix"
	URLBase    = "https://tnaflix.com"
	cookieFile = ".cookies_tnaflix.json"
)

var client = core.NewClient(cookieFile, URLBase)

// videoCard extracts one <div> from a video-list container into VideoData,
// shared between search results and embed suggestions (both use the same markup).
func videoCard(video *goquery.Selection, pageNumber int, urlFont string, indice int) core.VideoData {
	thumbA := video.Find("a.video-thumb").First()
	titleA := video.Find("a.video-title").First()

	videoURL, _ := thumbA.Attr("href")
	title := strings.TrimSpace(titleA.Text())
	timeVal := strings.TrimSpace(thumbA.Find("div.video-duration").First().Text())
	img := thumbA.Find("img").First()
	thumbnail, ok := img.Attr("data-src")
	if !ok {
		thumbnail, _ = img.Attr("src")
	}
	preview, _ := thumbA.Attr("data-trailer")

	eyeIcon := video.Find("i.icon-eye").First()
	var views string
	if eyeIcon.Length() > 0 {
		views = strings.TrimSpace(eyeIcon.Parent().Text())
	}
	thumbUpIcon := video.Find("i.icon-thumb-up").First()
	var rating string
	if thumbUpIcon.Length() > 0 {
		rating = strings.TrimSpace(thumbUpIcon.Parent().Text())
	}

	dur, _ := core.ParseDurationClock(timeVal)

	return core.VideoData{
		URL:             videoURL,
		Title:           title,
		Time:            timeVal,
		Thumbnail:       thumbnail,
		Rating:          rating,
		Preview:         preview,
		DurationSeconds: dur,
		Duration:        dur,
		URLFont:         urlFont,
		Views:           views,
		PageNumber:      pageNumber,
		SiteName:        SiteName,
		Indice:          indice,
	}
}

func minerateVideoList(doc *goquery.Document, pageNumber int, urlFont string) []core.VideoData {
	var videos []core.VideoData
	i := 0
	doc.Find("div.video-list").First().Children().Each(func(_ int, video *goquery.Selection) {
		videos = append(videos, videoCard(video, pageNumber, urlFont, i))
		i++
	})
	return videos
}

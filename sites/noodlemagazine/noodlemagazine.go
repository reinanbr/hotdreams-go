// Package noodlemagazine ports hotdreams.sites.noodlemagazine to Go.
package noodlemagazine

import (
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

const (
	SiteName   = "NoodleMagazine"
	URLBase    = "https://noodlemagazine.com"
	cookieFile = ".cookies_noodlemagazine.json"
)

var client = core.NewClient(cookieFile, URLBase)

// clockToSeconds mirrors utils.py's bespoke "H:MM:SS"/"MM:SS" parsing
// (kept separate from core.ParseDurationClock: the Python code computes
// minutes-carrying-hours before seconds, same result but written by hand).
func clockToSeconds(timeVideo string) (int, bool) {
	parts := strings.Split(timeVideo, ":")
	if len(parts) > 2 {
		hour, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		min, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		sec, err3 := strconv.Atoi(strings.TrimSpace(parts[2]))
		if err1 != nil || err2 != nil || err3 != nil {
			return 0, false
		}
		min = min + 60*hour
		return min*60 + sec, true
	}
	if len(parts) == 2 {
		min, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		sec, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err1 != nil || err2 != nil {
			return 0, false
		}
		return min*60 + sec, true
	}
	return 0, false
}

// minerateVideo mirrors noodlemagazine/utils.py's minerate_video. urlFont is
// the request URL passed in by the caller (not the post-redirect URL), same
// as the Python code.
func minerateVideo(doc *goquery.Document, urlFont string, pageNumber int) ([]core.VideoData, bool) {
	if strings.Contains(doc.Text(), "Nothing Found") {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div.item").Each(func(_ int, video *goquery.Selection) {
		href, _ := video.Find("a").First().Attr("href")
		img := video.Find("img").First()
		title, _ := img.Attr("alt")
		timeVideo := video.Find("div.m_time").First().Text()
		thumbnail, ok := img.Attr("data-src")
		if !ok {
			thumbnail, _ = img.Attr("src")
		}
		views := video.Find("div.m_views").First().Text()
		var preview string
		if iImg := video.Find("div.i_img").First(); iImg.Length() > 0 {
			preview, _ = iImg.Attr("data-trailer_url")
		}

		dur, _ := clockToSeconds(timeVideo)

		videos = append(videos, core.VideoData{
			Title:           title,
			Time:            timeVideo,
			DurationSeconds: dur,
			PageNumber:      pageNumber,
			URL:             URLBase + href,
			URLFont:         urlFont,
			Views:           views,
			Thumbnail:       thumbnail,
			SiteName:        SiteName,
			Preview:         preview,
			Indice:          i,
		})
		i++
	})
	return videos, true
}

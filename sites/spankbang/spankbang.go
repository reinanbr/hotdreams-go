// Package spankbang ports hotdreams.sites.spankbang to Go.
//
// NOTE: spankbang.com is behind Cloudflare's JS challenge. The Python lib
// works around it with the `cloudscraper` library; this port uses a plain
// HTTP client like every other site here, so — same as the Python README
// already flags this site as "blocked by Cloudflare" — requests may be
// blocked without a dedicated Cloudflare-challenge solver in front of it.
package spankbang

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "SpankBang"
	URLBase    = "https://spankbang.com"
	cookieFile = ".cookies_spankbang.json"
)

var client = core.NewClient(cookieFile, URLBase)

// stats holds the parsed and raw forms of a video-item's "stats" blob
// (views / rating / relative upload date), mirroring schoolbag.py's work_stats.
type stats struct {
	Views            int64
	Rating           int
	DateUploadSecond float64
	ViewsStr         string
	RatingStr        string
	DateUploadStr    string
}

// workStats mirrors schoolbag.py's work_stats: it expects the stats blob to
// split (on "\n") into ["", views, rating, relative-date, ...], same
// assumption the Python code makes about the site's markup.
func workStats(raw string) stats {
	raw = strings.ReplaceAll(raw, "\n\n\n", "")
	parts := strings.Split(raw, "\n")
	get := func(i int) string {
		if i < len(parts) {
			return parts[i]
		}
		return ""
	}

	viewsStr := get(1)
	var views int64
	switch {
	case strings.Contains(viewsStr, "M"):
		if f, err := strconv.ParseFloat(strings.Split(viewsStr, "M")[0], 64); err == nil {
			views = int64(f * 1e6)
		}
	case strings.Contains(viewsStr, "K"):
		if f, err := strconv.ParseFloat(strings.Split(viewsStr, "K")[0], 64); err == nil {
			views = int64(f * 1e3)
		}
	default:
		if f, err := strconv.ParseFloat(viewsStr, 64); err == nil {
			views = int64(f)
		}
	}

	ratingStr := get(2)
	rating := core.AtoiOr(strings.Split(ratingStr, "%")[0], 0)

	dateUploadStr := get(3)
	const hourSeconds = 3600.0
	now := float64(time.Now().Unix())
	var dateUpload float64
	switch {
	case strings.Contains(dateUploadStr, "minute"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " minute")[0], 0)
		dateUpload = now - float64(n)*60
	case strings.Contains(dateUploadStr, "hour"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " hour")[0], 0)
		dateUpload = now - float64(n)*hourSeconds
	case strings.Contains(dateUploadStr, "day"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " day")[0], 0)
		dateUpload = now - float64(n)*24*hourSeconds
	case strings.Contains(dateUploadStr, "week"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " week")[0], 0)
		dateUpload = now - float64(n)*7*24*hourSeconds
	case strings.Contains(dateUploadStr, "month"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " month")[0], 0)
		dateUpload = now - float64(n)*30*24*hourSeconds
	case strings.Contains(dateUploadStr, "year"):
		n := core.AtoiOr(strings.Split(dateUploadStr, " year")[0], 0)
		dateUpload = now - float64(n)*365*24*hourSeconds
	default:
		dateUpload = math.NaN()
	}

	return stats{
		Views:            views,
		Rating:           rating,
		DateUploadSecond: dateUpload,
		ViewsStr:         viewsStr,
		RatingStr:        ratingStr,
		DateUploadStr:    dateUploadStr,
	}
}

func durationFromMinutesLabel(timeVal string) int {
	mins := core.AtoiOr(strings.Split(timeVal, " min")[0], 0)
	return mins * 60
}

// minerateSearchResults mirrors search_porn.py's get_videos_bg_link parsing
// of div.video-item cards on a search-results page.
func minerateSearchResults(doc *goquery.Document, urlFont string, pageNumber int) ([]core.VideoData, bool) {
	videoDiv := doc.Find("div.video-item")
	if videoDiv.Length() == 0 || strings.Contains(doc.Text(), "We could not find any videos for ") {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	videoDiv.Each(func(_ int, video *goquery.Selection) {
		href, _ := video.Find("a").First().Attr("href")
		img := video.Find("img").First()
		title, _ := img.Attr("alt")
		timeVal := video.Find("span.l").First().Text()
		thumbnail, _ := img.Attr("data-src")
		rawStats := video.Find("div.stats").First().Text()
		rawStats = strings.ReplaceAll(strings.ReplaceAll(rawStats, "aa", " "), "\n\n\n", "")
		sw := workStats(rawStats)
		gifURL, _ := img.Attr("data-preview")

		videos = append(videos, core.VideoData{
			Title:             title,
			Time:              timeVal,
			DurationSeconds:   durationFromMinutesLabel(timeVal),
			Views:             sw.ViewsStr,
			Rating:            sw.RatingStr,
			DateUpload:        sw.DateUploadStr,
			ViewsInt:          sw.Views,
			DateUploadSeconds: sw.DateUploadSecond,
			RatingInt:         sw.Rating,
			PageNumber:        pageNumber,
			URL:               URLBase + href,
			URLFont:           urlFont,
			Thumbnail:         thumbnail,
			Preview:           gifURL,
			SiteName:          SiteName,
			Indice:            i,
		})
		i++
	})
	return videos, true
}

// minerateSuggestions mirrors embed.py's suggestion-list parsing of the same
// div.video-item cards, but keeps the stats blob as raw text (no work_stats
// parsing) — same distinction the Python code makes.
func minerateSuggestions(doc *goquery.Document, urlFont string) []core.VideoData {
	var videos []core.VideoData
	i := 0
	doc.Find("div.video-item").Each(func(_ int, video *goquery.Selection) {
		href, _ := video.Find("a").First().Attr("href")
		img := video.Find("img").First()
		title, _ := img.Attr("alt")
		timeVal := video.Find("span.l").First().Text()
		thumbnail, _ := img.Attr("data-src")

		var statsText string
		if s := video.Find("div.stats").First(); s.Length() > 0 {
			statsText = s.Text()
			replacer := strings.NewReplacer(
				"\n", " ", "\n\n", "", "(", "", ")", "", "'", "", " ", "", ",", "",
			)
			statsText = replacer.Replace(statsText)
		}
		gifURL, _ := img.Attr("data-preview")

		videos = append(videos, core.VideoData{
			Title:      title,
			Time:       timeVal,
			Duration:   durationFromMinutesLabel(timeVal),
			Stats:      statsText,
			PageNumber: 0,
			URL:        URLBase + href,
			URLFont:    urlFont,
			Thumbnail:  thumbnail,
			Preview:    gifURL,
			SiteName:   SiteName,
			Indice:     i,
		})
		i++
	})
	return videos
}

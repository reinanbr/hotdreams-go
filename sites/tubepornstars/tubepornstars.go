// Package tubepornstars scrapes tubepornstars.com, a pornstar-oriented video
// aggregator: every result links out to the tube that hosts the video
// (Eporner, xHamster, PornOne...), so there's no embed page of its own.
package tubepornstars

import (
	"encoding/base64"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "tubepornstars"
	URLBase    = "https://www.tubepornstars.com"
	cookieFile = ".cookies_tubepornstars.json"

	noResultsText = "No results found for the following search"
)

var client = core.NewClient(cookieFile, URLBase)

// sourceURL extracts the original video URL from a "/out/?l=..." redirect
// link. The l param is base64 of a msgpack array whose first http(s) string
// is the target URL. Returns "" if it can't be decoded.
func sourceURL(outLink string) string {
	u, err := url.Parse(outLink)
	if err != nil {
		return ""
	}
	l := strings.TrimRight(u.Query().Get("l"), "=")
	data, err := base64.RawStdEncoding.DecodeString(l)
	if err != nil {
		return ""
	}
	for i := 0; i < len(data); i++ {
		var n, start int
		switch b := data[i]; {
		case b >= 0xa0 && b <= 0xbf: // fixstr
			n, start = int(b&0x1f), i+1
		case b == 0xd9 && i+1 < len(data): // str8
			n, start = int(data[i+1]), i+2
		case b == 0xda && i+2 < len(data): // str16
			n, start = int(data[i+1])<<8|int(data[i+2]), i+3
		default:
			continue
		}
		if start+n > len(data) {
			continue
		}
		if s := string(data[start : start+n]); strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
			return s
		}
	}
	return ""
}

// minerateVideo parses the video cards of a search/pornstar page, skipping
// the paid (ad) ones.
func minerateVideo(doc *goquery.Document, pageNumber int, urlFont string) ([]core.VideoData, bool) {
	if strings.Contains(doc.Find("h3").Text(), noResultsText) {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div.card.sub").Each(func(_ int, card *goquery.Selection) {
		if card.HasClass("paid") {
			return
		}
		link := card.Find("a.item-title").First()
		href, _ := link.Attr("href")
		if href == "" {
			return
		}
		outURL := core.ResolveURL(URLBase, href)
		videoURL := sourceURL(outURL)
		if videoURL == "" {
			videoURL = outURL
		}
		title, _ := link.Attr("title")
		if title == "" {
			title = strings.TrimSpace(link.Text())
		}

		// badge reads e.g. "HD 31:06": the clock is the last field
		var timeVal string
		if f := strings.Fields(card.Find("span.badge").First().Text()); len(f) > 0 {
			timeVal = f[len(f)-1]
		}
		dur, _ := core.ParseDurationClock(timeVal)

		rating := strings.TrimSpace(card.Find("span.item-score").First().Text())
		meta := card.Find("div.item-source-rating-container")
		source := strings.TrimSpace(meta.Find("a.item-source").First().Text())
		uploaded := strings.TrimSpace(meta.Find("span.item-source").First().Text())

		videos = append(videos, core.VideoData{
			Title:      title,
			Time:       timeVal,
			Duration:   dur,
			PageNumber: pageNumber,
			URL:        videoURL,
			URLFont:    urlFont,
			Thumbnail:  core.ImageURL(card.Find("img.item-image").First(), URLBase),
			Stats:      source,
			Rating:     rating,
			RatingInt:  core.AtoiOr(strings.TrimSuffix(rating, "%"), 0),
			DateUpload: uploaded,
			SiteName:   SiteName,
			Indice:     i,
		})
		i++
	})
	return videos, true
}

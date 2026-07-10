// Package nudes7 ports hotdreams.sites.only.nudes7 to Go.
//
// nudes7.com now permanently redirects to shareanynudes.com (same content,
// new domain and markup); this port follows it automatically, same as the
// Python lib.
package nudes7

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "nudes7@only"
	URLBase    = "https://shareanynudes.com/"
	cookieFile = ".cookies_nudes7.json"
)

var client = core.NewClient(cookieFile, URLBase)

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	var body []byte
	var finalURL string
	var lastErr error
	for i := 0; i < 6; i++ {
		var err error
		body, finalURL, _, err = client.Fetch(rawURL)
		if err == nil {
			lastErr = nil
			break
		}
		lastErr = err
		time.Sleep(time.Second)
	}
	if lastErr != nil {
		return nil, false, lastErr
	}

	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div.item").Each(func(_ int, item *goquery.Selection) {
		if item.HasClass("ad-item") {
			return
		}
		link := item.Find("a").First()
		href, hasHref := link.Attr("href")
		if link.Length() == 0 || !hasHref {
			return
		}
		title, _ := link.Attr("title")
		var thumbnail string
		if img := item.Find("img").First(); img.Length() > 0 {
			thumbnail, _ = img.Attr("src")
		}
		duration := item.Find("div.duration").First()
		views := item.Find("div.views").First()

		videos = append(videos, core.VideoData{
			Title:      title,
			URL:        href,
			Thumbnail:  thumbnail,
			Time:       strings.TrimSpace(duration.Text()),
			Views:      strings.TrimSpace(views.Text()),
			SiteName:   SiteName,
			URLFont:    finalURL,
			Indice:     i,
			PageNumber: pageNumber,
		})
		i++
	})
	return videos, true, nil
}

func buildSearchURL(query string, page int) string {
	return "https://shareanynudes.com/search/" + query + "/"
}

// SearchPorn searches nudes7 (shareanynudes.com) for query, a single page.
func SearchPorn(query string) (core.DataVideos, error) {
	pageNumber := 1
	return core.SearchPornBase(query, URLBase, SiteName, 0, &pageNumber, fetchPage, buildSearchURL)
}

// Package viralpornhub ports hotdreams.sites.only.viralpornhub to Go.
package viralpornhub

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "viralpornhub@only"
	URLBase    = "https://viralpornhub.com"
	cookieFile = ".cookies_viralpornhub.json"
)

var client = core.NewClient(cookieFile, URLBase)

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, finalURL, _, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div.item").Each(func(_ int, item *goquery.Selection) {
		a := item.Find("a").First()
		title, _ := a.Attr("title")
		href, _ := a.Attr("href")

		var thumbnail string
		if dataWebp, ok := item.Find("img").First().Attr("data-webp"); ok {
			if parts := strings.SplitN(dataWebp, "/ss.", 2); len(parts) == 2 {
				thumbnail = "https://" + parts[1]
			}
		}
		preview, _ := item.Find("img").First().Attr("data-preview")
		timeVal := item.Find("div.duration").First().Text()
		rating := strings.Join(strings.Fields(item.Find("div.rating").First().Text()), "")
		dateUpload := item.Find("div.added").First().Text()
		views := item.Find("div.views").First().Text()

		videos = append(videos, core.VideoData{
			Title:      title,
			URL:        href,
			Thumbnail:  thumbnail,
			Preview:    preview,
			Rating:     rating,
			Time:       timeVal,
			DateUpload: dateUpload,
			Views:      views,
			PageNumber: 1,
			Indice:     i,
			URLFont:    finalURL,
			SiteName:   SiteName,
		})
		i++
	})
	return videos, true, nil
}

func buildSearchURL(query string, page int) string {
	return "https://viralpornhub.com/search/" + query
}

// SearchPorn searches viralpornhub for query. Mirrors the Python lib's call
// (page_limit=1, page_number=None): a single paginated iteration, not the
// direct single-page code path.
func SearchPorn(query string) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, 1, nil, fetchPage, buildSearchURL)
}

// Package ukdevilz ports hotdreams.sites.ukdevilz to Go.
package ukdevilz

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

const (
	SiteName   = "ukdevilz"
	URLBase    = "https://ukdevilz.com"
	cookieFile = ".cookies_ukdevilz.json"
)

var client = core.NewClient(cookieFile, URLBase)

func minerateVideo(doc *goquery.Document, urlFont string, pageNumber int) ([]core.VideoData, bool) {
	if strings.Contains(doc.Text(), "Nothing Found") {
		return nil, false
	}

	var videos []core.VideoData
	i := 0
	doc.Find("div.item").Each(func(_ int, video *goquery.Selection) {
		href, _ := video.Find("a").First().Attr("href")
		title, _ := video.Find("img").First().Attr("alt")
		timeVideo := video.Find("div.m_time").First().Text()
		thumbnail, _ := video.Find("img").First().Attr("data-src")
		views := video.Find("div.m_views").First().Text()

		var dur int
		parts := strings.SplitN(timeVideo, ":", 2)
		if len(parts) == 2 {
			if mins, err1 := strconv.Atoi(strings.TrimSpace(parts[0])); err1 == nil {
				if secs, err2 := strconv.Atoi(strings.TrimSpace(parts[1])); err2 == nil {
					dur = mins*60 + secs
				}
			}
		}

		videos = append(videos, core.VideoData{
			Title:      title,
			Time:       timeVideo,
			Duration:   dur,
			PageNumber: pageNumber,
			URL:        URLBase + href,
			URLFont:    urlFont,
			Stats:      views,
			Thumbnail:  thumbnail,
			SiteName:   SiteName,
			Indice:     i,
		})
		i++
	})
	return videos, true
}

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/video/%s?p=%d", URLBase, query, page)
}

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, _, _, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}
	videos, ok := minerateVideo(doc, rawURL, pageNumber)
	return videos, ok, nil
}

// SearchPorn searches ukdevilz for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

// GetVideoEmbed fetches a single ukdevilz video's embed url (that's the only
// data the site's embed page exposes).
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	var link string
	if iframe := doc.Find("iframe").First(); iframe.Length() > 0 {
		if src, ok := iframe.Attr("src"); ok {
			link = URLBase + strings.Split(src, "&a=1")[0]
		}
	}

	return core.EmbedVideo{
		SiteName:            SiteName,
		URL:                 link,
		LenVideosSugestions: 0,
	}, nil
}

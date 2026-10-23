// Package pornone ports hotdreams.sites.pornone to Go.
package pornone

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hotdreams-go/core"
)

const (
	SiteName   = "pornone"
	URLBase    = "https://pornone.com"
	cookieFile = ".cookies_pornone.json"
)

var client = core.NewClient(cookieFile, URLBase)

// normalizeStats approximates bs4's `.get_text(' ', strip=True)` followed by
// the Python code's newline cleanup: collapse all whitespace to single spaces.
func normalizeStats(sel *goquery.Selection) string {
	return strings.Join(strings.Fields(sel.Text()), " ")
}

func videoCards(heaven *goquery.Selection, pageNumber int, urlFont string) []core.VideoData {
	var videos []core.VideoData
	i := 0
	heaven.Find("a.videocard").Each(func(_ int, vd *goquery.Selection) {
		titleDiv := vd.Find("div.videotitle").First()
		if titleDiv.Length() == 0 {
			return
		}
		title := titleDiv.Text()
		timeVal := strings.TrimSpace(vd.Find("span.durlabel").First().Text())
		stats := normalizeStats(vd.Find("div.author").First())
		img := vd.Find("img.imgvideo").First()
		if img.Length() == 0 {
			img = vd.Find("img").First()
		}
		urlImg := core.ImageURL(img, URLBase)
		urlVideo, _ := vd.Attr("href")
		if urlVideo != "" {
			urlVideo = core.ResolveURL(URLBase, urlVideo)
		}

		var dur int
		parts := strings.SplitN(timeVal, ":", 2)
		if len(parts) == 2 {
			if mins, err1 := strconv.Atoi(strings.TrimSpace(parts[0])); err1 == nil {
				if secs, err2 := strconv.Atoi(strings.TrimSpace(parts[1])); err2 == nil {
					dur = mins*60 + secs
				}
			}
		}

		videos = append(videos, core.VideoData{
			Title:      title,
			Time:       strings.ReplaceAll(timeVal, "\n", ""),
			Duration:   dur,
			PageNumber: pageNumber,
			URL:        urlVideo,
			URLFont:    urlFont,
			Stats:      stats,
			Thumbnail:  urlImg,
			SiteName:   SiteName,
			Indice:     i,
		})
		i++
	})
	return videos
}

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, finalURL, _, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	rawText := string(body)
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}

	if strings.Contains(rawText, "Pornstar rank") {
		if strings.Contains(doc.Text(), "No videos for this user.") {
			return nil, false, nil
		}
		urlFont := fmt.Sprintf("%s%d/", finalURL, pageNumber)
		heaven := doc.Find("div").FilterFunction(func(_ int, s *goquery.Selection) bool {
			return s.HasClass("portrait:gap-x-0.1")
		}).First()
		if heaven.Length() == 0 {
			return nil, false, nil
		}
		return videoCards(heaven, pageNumber, urlFont), true, nil
	}

	if strings.Contains(doc.Text(), "No results found for this criteria.") {
		return nil, false, nil
	}

	checkQuery := query
	if strings.Contains(checkQuery, " ") {
		fields := strings.Split(checkQuery, " ")
		cleaned := fields[:0]
		for _, f := range fields {
			if f != "" {
				cleaned = append(cleaned, f)
			}
		}
		checkQuery = strings.Join(cleaned, " ")
	}
	lowerText := strings.ToLower(doc.Text())
	if strings.Contains(checkQuery, " ") {
		first := strings.SplitN(checkQuery, " ", 2)[0]
		if strings.Count(lowerText, strings.ToLower(first)) <= 2 {
			return nil, false, nil
		}
	} else if strings.Count(lowerText, strings.ToLower(checkQuery)) <= 2 {
		return nil, false, nil
	}

	heaven := doc.Find("div.mt-1").First()
	if heaven.Length() == 0 {
		return nil, false, nil
	}
	return videoCards(heaven, pageNumber, finalURL), true, nil
}

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/search/?q=%s&sort=relevance&filter=&page=%d", URLBase, query, page)
}

// SearchPorn searches pornone for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

// GetVideoEmbed fetches a single pornone video's embed data.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	var link, thumbnail string
	if video := doc.Find("video").First(); video.Length() > 0 {
		link, _ = video.Find("source").First().Attr("src")
		if poster, ok := video.Attr("poster"); ok && poster != "" {
			thumbnail = core.ResolveURL(URLBase, poster)
		}
	}
	if thumbnail == "" {
		if og := core.Meta(doc, "property", "og:image"); og != "" {
			thumbnail = core.ResolveURL(URLBase, og)
		}
	}
	title := core.Text(doc, "h1.text-3xl")

	return core.EmbedVideo{
		SiteName:            SiteName,
		URL:                 link,
		Title:               title,
		Thumbnail:           thumbnail,
		LenVideosSugestions: 0,
	}, nil
}

package eporner

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(s), 10, 64)
}

func stringify(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		if t == float64(int64(t)) {
			return strconv.FormatInt(int64(t), 10)
		}
		return fmt.Sprintf("%v", t)
	default:
		return fmt.Sprintf("%v", t)
	}
}

// videoObject finds the <script type="application/ld+json"> whose @type is
// VideoObject, mirroring embed.py's _video_object.
func videoObject(doc *goquery.Document) map[string]interface{} {
	var result map[string]interface{}
	doc.Find(`script[type="application/ld+json"]`).EachWithBreak(func(_ int, s *goquery.Selection) bool {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(s.Text()), &data); err != nil {
			return true
		}
		if t, _ := data["@type"].(string); t == "VideoObject" {
			result = data
			return false
		}
		return true
	})
	return result
}

func stringOrFirst(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case []interface{}:
		if len(t) > 0 {
			if s, ok := t[0].(string); ok {
				return s
			}
		}
	}
	return ""
}

// GetVideoEmbed fetches a single Eporner video's embed data and suggestions.
// Its returned URL is already the direct, playable .mp4 file.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	videoInfo := videoObject(doc)
	if videoInfo == nil {
		videoInfo = map[string]interface{}{}
	}

	title, _ := videoInfo["name"].(string)
	link, _ := videoInfo["contentUrl"].(string)
	if link == "" {
		link, _ = videoInfo["embedUrl"].(string)
	}
	uploadDate, _ := videoInfo["uploadDate"].(string)

	thumbnail := stringOrFirst(videoInfo["image"])
	if thumbnail == "" {
		thumbnail = core.Meta(doc, "property", "og:image")
	}

	durationSeconds := core.AtoiOr(core.Meta(doc, "property", "og:duration"), 0)

	var viewsInt int64
	if interaction, ok := videoInfo["interactionStatistic"].(map[string]interface{}); ok {
		switch n := interaction["userInteractionCount"].(type) {
		case float64:
			viewsInt = int64(n)
		case string:
			if v, err := parseInt64(n); err == nil {
				viewsInt = v
			}
		}
	}
	views := core.Text(doc, "div#cinemaviews1")

	var rating string
	if agg, ok := videoInfo["aggregateRating"].(map[string]interface{}); ok {
		if rv, ok := agg["ratingValue"]; ok && rv != nil {
			rating = stringify(rv) + "%"
		}
	}

	var actors []map[string]interface{}
	switch a := videoInfo["actor"].(type) {
	case map[string]interface{}:
		actors = []map[string]interface{}{a}
	case []interface{}:
		for _, item := range a {
			if m, ok := item.(map[string]interface{}); ok {
				actors = append(actors, m)
			}
		}
	}
	var names []string
	for _, a := range actors {
		if n, _ := a["name"].(string); n != "" {
			names = append(names, n)
		}
	}
	person := strings.Join(names, ", ")

	var tags []string
	doc.Find("div#video-info-tags a").Each(func(_ int, a *goquery.Selection) {
		tags = append(tags, strings.TrimSpace(a.Text()))
	})

	suggestions, _ := minerateVideo(doc, url, 0)

	return core.EmbedVideo{
		Title:               title,
		URL:                 link,
		DurationSeconds:     durationSeconds,
		UploadDate:          uploadDate,
		Views:               views,
		ViewsInt:            viewsInt,
		Rating:              rating,
		Thumbnail:           thumbnail,
		Person:              person,
		Tags:                strings.Join(tags, ", "),
		SiteName:            SiteName,
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

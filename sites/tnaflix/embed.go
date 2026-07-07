package tnaflix

import (
	"encoding/json"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/reinanbr/hot_dreams/core"
)

func firstOfAny(v interface{}) string {
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

// GetVideoEmbed fetches a single tnaflix video's embed data and suggestions.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	if strings.Contains(url, "?") {
		url += "&autoPlay=1"
	} else {
		url += "?autoPlay=1"
	}

	body, finalURL, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	var videoInfo map[string]interface{}
	scriptText := doc.Find(`script[type="application/ld+json"]`).First().Text()
	_ = json.Unmarshal([]byte(scriptText), &videoInfo)

	title, _ := videoInfo["name"].(string)
	thumbnail := firstOfAny(videoInfo["thumbnailUrl"])
	link, _ := videoInfo["embedUrl"].(string)

	badges := doc.Find("div.video-detail-badges").First()
	var views string
	eyeIcon := badges.Find("i.icon-eye").First()
	if eyeIcon.Length() > 0 {
		views = strings.TrimSpace(eyeIcon.Parent().Text())
	}
	var tagParts []string
	badges.Find("a.badge-video").Each(func(_ int, a *goquery.Selection) {
		tagParts = append(tagParts, strings.TrimSpace(a.Text()))
	})

	suggestions := minerateVideoList(doc, 0, finalURL)

	return core.EmbedVideo{
		SiteName:            SiteName,
		URL:                 link,
		Title:               title,
		Thumbnail:           thumbnail,
		Views:               views,
		Tags:                strings.Join(tagParts, ", "),
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

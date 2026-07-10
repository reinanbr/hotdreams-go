package fpo

import (
	"regexp"

	"github.com/reinanbr/hotdreams-go/core"
)

var videoCategoriesRe = regexp.MustCompile(`video_categories:\s*'([^']*)'`)

// GetVideoEmbed fetches a single FPO.XXX video's embed data and suggestions.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	title := core.Meta(doc, "property", "og:title")
	link := core.Meta(doc, "property", "og:video")
	if link == "" {
		link = core.Meta(doc, "property", "og:video:url")
	}
	thumbnail := core.Meta(doc, "property", "og:image")
	uploadDate := core.Meta(doc, "property", "ya:ovs:upload_date")
	durationSeconds := core.AtoiOr(core.Meta(doc, "property", "video:duration"), 0)
	viewsInt := int64(core.AtoiOr(core.Meta(doc, "property", "ya:ovs:views_total"), 0))
	rating := core.Text(doc, "span.voters")
	person := core.Meta(doc, "property", "ya:ovs:login")

	var tags string
	if m := videoCategoriesRe.FindStringSubmatch(string(body)); len(m) == 2 && m[1] != "" {
		tags = m[1]
	}

	suggestions, _ := minerateVideo(doc, url, 0)

	return core.EmbedVideo{
		Title:               title,
		URL:                 link,
		DurationSeconds:     durationSeconds,
		UploadDate:          uploadDate,
		ViewsInt:            viewsInt,
		Rating:              rating,
		Thumbnail:           thumbnail,
		Person:              person,
		Tags:                tags,
		SiteName:            SiteName,
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

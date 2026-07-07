package noodlemagazine

import (
	"github.com/reinanbr/hot_dreams/core"
)

// GetVideoEmbed fetches a single NoodleMagazine video's embed data and suggestions.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	suggestions, _ := minerateVideo(doc, url, 0)

	return core.EmbedVideo{
		Title:               core.Meta(doc, "property", "og:title"),
		URL:                 core.Meta(doc, "property", "og:video"),
		DurationSeconds:     core.AtoiOr(core.Meta(doc, "property", "video:duration"), 0),
		UploadDate:          core.Meta(doc, "property", "ya:ovs:upload_date"),
		ViewsInt:            int64(core.AtoiOr(core.Meta(doc, "property", "ya:ovs:views_total"), 0)),
		Likes:               core.Meta(doc, "property", "ya:ovs:likes"),
		Thumbnail:           core.Meta(doc, "property", "og:image"),
		Person:              core.Meta(doc, "property", "ya:ovs:person"),
		Tags:                core.Meta(doc, "property", "video:tag"),
		SiteName:            SiteName,
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

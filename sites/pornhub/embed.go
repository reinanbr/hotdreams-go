package pornhub

import (
	"github.com/reinanbr/hot_dreams/core"
)

// GetVideoEmbed fetches a single PornHub video's embed data and suggestions.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, finalURL, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	suggestions, _ := minerateVideo(doc, 0, finalURL)

	return core.EmbedVideo{
		Title:               core.Meta(doc, "property", "og:title"),
		URL:                 core.Meta(doc, "name", "twitter:player"),
		DurationSeconds:     core.AtoiOr(core.Meta(doc, "property", "video:duration"), 0),
		Thumbnail:           core.Attr(doc, `link[as="image"]`, "href"),
		Likes:               core.Text(doc, "div.votes-fav-wrap"),
		UploadDate:          core.Text(doc, "div.videoInfo"),
		Views:               core.Text(doc, "div.views"),
		Rating:              core.Text(doc, "div.ratingPercent"),
		Tags:                core.Text(doc, "div.tagsWrapper"),
		Person:              core.Text(doc, "div.pornstarsWrapper"),
		SiteName:            SiteName,
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

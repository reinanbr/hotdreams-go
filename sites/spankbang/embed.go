package spankbang

import (
	"github.com/reinanbr/hotdreams-go/core"
)

// GetVideoEmbed fetches a single SpankBang video's embed data and suggestions.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	title := core.Text(doc, "h1")
	timeVal := core.Text(doc, "span.i-length")
	players := core.Text(doc, "span.i-plays")
	timePublished := core.Text(doc, "span.i-date")
	thumbnail := core.Attr(doc, "img.player_thumb", "src")
	link := core.Attr(doc, "video source", "src")

	suggestions := minerateSuggestions(doc, url)

	return core.EmbedVideo{
		SiteName:            SiteName,
		URL:                 link,
		Title:               title,
		Time:                timeVal,
		Thumbnail:           thumbnail,
		TimePublished:       timePublished,
		Views:               players,
		LenVideosSugestions: len(suggestions),
		VideosSugestions:    suggestions,
	}, nil
}

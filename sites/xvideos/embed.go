package xvideos

import (
	"fmt"
	"strings"

	"github.com/reinanbr/hotdreams-go/core"
)

// GetVideoEmbed fetches a single XVideos video's embed data.
func GetVideoEmbed(url string) (core.EmbedVideo, error) {
	body, _, _, err := client.Fetch(url)
	if err != nil {
		return core.EmbedVideo{}, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return core.EmbedVideo{}, err
	}

	segments := strings.Split(url, "video")
	if len(segments) < 3 {
		return core.EmbedVideo{}, fmt.Errorf("xvideos: could not extract video key from url %q", url)
	}
	keyVideo := strings.Split(segments[2], "/")[0]

	return core.EmbedVideo{
		Title:               core.Meta(doc, "property", "og:title"),
		URL:                 fmt.Sprintf("%s/embedframe/%s", URLBase, keyVideo),
		DurationSeconds:     core.AtoiOr(core.Meta(doc, "property", "og:duration"), 0),
		Thumbnail:           core.Meta(doc, "property", "og:image"),
		Views:               core.Text(doc, "strong.mobile-hide"),
		Rating:              core.Text(doc, "div.rating-good-perc.hide-if-zero-2045"),
		SiteName:            SiteName,
		LenVideosSugestions: 0,
	}, nil
}

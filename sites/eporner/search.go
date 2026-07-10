package eporner

import (
	"fmt"

	"github.com/reinanbr/hotdreams-go/core"
)

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/search/%s/%d/", URLBase, query, page)
}

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, _, status, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	if status == 404 {
		return nil, false, nil
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}
	videos, ok := minerateVideo(doc, rawURL, pageNumber)
	return videos, ok, nil
}

// SearchPorn searches Eporner for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

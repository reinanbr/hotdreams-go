package spankbang

import (
	"fmt"

	"github.com/reinanbr/hot_dreams/core"
)

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/s/%s/%d/?o=trending", URLBase, query, page)
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
	videos, ok := minerateSearchResults(doc, rawURL, pageNumber)
	return videos, ok, nil
}

// SearchPorn searches SpankBang for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

package fpo

import (
	"fmt"

	"github.com/reinanbr/hotdreams-go/core"
)

// buildSearchURL mirrors utils.py's url_base_page_number_search: page 1 is a
// normal HTML page, later pages go through the same ajax endpoint the site's
// own "load more" pagination uses.
func buildSearchURL(query string, page int) string {
	urlSearch := fmt.Sprintf("%s/search/%s/", URLBase, query)
	if page <= 1 {
		return urlSearch
	}
	return fmt.Sprintf(
		"%s?mode=async&function=get_block&block_id=list_videos_videos_list_search_result&q=%s&category_ids=&sort_by=&from_videos=%d&from_albums=%d",
		urlSearch, query, page, page,
	)
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
	videos, ok := minerateVideo(doc, rawURL, pageNumber)
	return videos, ok, nil
}

// SearchPorn searches FPO.XXX for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

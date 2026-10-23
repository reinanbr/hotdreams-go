package tubepornstars

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/reinanbr/hotdreams-go/core"
)

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/search/a/%s?page=%d", URLBase, query, page)
}

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, finalURL, status, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	// a pornstar-name query redirects to /pornstar/<slug>, dropping ?page=N
	if pageNumber > 1 && status == http.StatusOK {
		if u, err := url.Parse(finalURL); err == nil && u.Query().Get("page") == "" {
			q := u.Query()
			q.Set("page", fmt.Sprint(pageNumber))
			u.RawQuery = q.Encode()
			body, finalURL, status, err = client.Fetch(u.String())
			if err != nil {
				return nil, false, err
			}
		}
	}
	// past the last page the site answers 404
	if status != http.StatusOK {
		return nil, false, nil
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}
	videos, ok := minerateVideo(doc, pageNumber, finalURL)
	return videos, ok, nil
}

// SearchPorn searches tubepornstars for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given. Each video's URL points to the
// tube hosting it.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

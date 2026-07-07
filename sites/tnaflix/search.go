package tnaflix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reinanbr/hot_dreams/core"
)

func buildSearchURL(query string, page int) string {
	return fmt.Sprintf("%s/search.php?tab=videos&what=%s&category=&sb=relevance&su=anytime&sd=all&dir=desc&page=%d", URLBase, query, page)
}

// actualPageFromURL extracts the &page= value tnaflix redirected us to,
// mirroring the Python code's `page_number > int(url_htm.url.split('&page=')[1])` guard.
func actualPageFromURL(finalURL string) (int, bool) {
	idx := strings.Index(finalURL, "&page=")
	if idx == -1 {
		return 0, false
	}
	rest := finalURL[idx+len("&page="):]
	if amp := strings.IndexByte(rest, '&'); amp != -1 {
		rest = rest[:amp]
	}
	n, err := strconv.Atoi(rest)
	if err != nil {
		return 0, false
	}
	return n, true
}

func fetchPage(rawURL string, pageNumber int, query string) ([]core.VideoData, bool, error) {
	body, finalURL, _, err := client.Fetch(rawURL)
	if err != nil {
		return nil, false, err
	}
	doc, err := core.ParseHTML(body)
	if err != nil {
		return nil, false, err
	}
	text := doc.Text()

	if strings.Contains(query, " ") {
		first := strings.SplitN(query, " ", 2)[0]
		if strings.Count(text, first) < 3 {
			return nil, false, nil
		}
	} else if strings.Count(text, query) < 3 {
		return nil, false, nil
	}

	if actualPage, ok := actualPageFromURL(finalURL); ok && pageNumber > actualPage {
		return nil, false, nil
	}

	if doc.Find("div.video-list").Length() == 0 {
		return nil, false, nil
	}

	return minerateVideoList(doc, pageNumber, finalURL), true, nil
}

// SearchPorn searches tnaflix for query, fetching up to pageLimit pages
// (default 2), or exactly pageNumber if given.
func SearchPorn(query string, pageLimit int, pageNumber *int) (core.DataVideos, error) {
	return core.SearchPornBase(query, URLBase, SiteName, pageLimit, pageNumber, fetchPage, buildSearchURL)
}

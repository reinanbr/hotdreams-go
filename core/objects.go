// Package core provides the shared plumbing (HTTP client, cookies, settings,
// data types) used by every site scraper in hot_dreams.
package core

import "fmt"

// VideoData mirrors hotdreams.utils.objects.VideoData from the Python library:
// a single search-result / suggestion video entry.
type VideoData struct {
	Title             string
	Time              string
	PageNumber        int
	URL               string
	URLFont           string
	Thumbnail         string
	SiteName          string
	Preview           string
	Stats             string
	Indice            int
	Views             string
	ViewsInt          int64
	Rating            string
	RatingInt         int
	DateUpload        string
	DateUploadSeconds float64
	Duration          int
	DurationSeconds   int
}

func (v VideoData) String() string {
	return fmt.Sprintf(
		"%s_video(title=%q, time=%q, url=%q, url_font=%q, thumbnail=%q, site_name=%q, page_number=%d, indice=%d)",
		v.SiteName, v.Title, v.Time, v.URL, v.URLFont, v.Thumbnail, v.SiteName, v.PageNumber, v.Indice,
	)
}

// DataVideos mirrors hotdreams.utils.objects.DataVideos: the result of a
// search_porn call. LenPages is nil when the caller asked for a single
// page (page_number set), mirroring Python's len_pages=None in that case.
type DataVideos struct {
	SiteName       string
	URLBase        string
	Query          string
	Ping           float64
	Videos         []VideoData
	VideosPerPages int
	LenPages       *int
	LenVideos      int
	URLSearch      string
	Success        bool
}

func (d DataVideos) String() string {
	lenPages := "None"
	if d.LenPages != nil {
		lenPages = fmt.Sprintf("%d", *d.LenPages)
	}
	return fmt.Sprintf(
		"%s(success=%v, query=%q, len_videos=%d, len_pages=%s, videos_per_pages=%d, ping=%v, url_base=%q, url_search=%q)",
		d.SiteName, d.Success, d.Query, d.LenVideos, lenPages, d.VideosPerPages, d.Ping, d.URLBase, d.URLSearch,
	)
}

// EmbedVideo mirrors hotdreams.utils.objects.EmbedVideo: the result of a
// get_video_embed call.
type EmbedVideo struct {
	URL                 string
	Title               string
	Time                string
	Thumbnail           string
	Views               string
	TimePublished       string
	SiteName            string
	LenVideosSugestions int
	VideosSugestions    []VideoData
	Rating              string
	DurationSeconds     int
	UploadDate          string
	ViewsInt            int64
	Likes               string
	Person              string
	Tags                string
}

func (e EmbedVideo) String() string {
	return fmt.Sprintf(
		"EmbedVideo%s(title=%q, url=%q, thumbnail=%q, views=%q, len_videos_sugestions=%d)",
		e.SiteName, e.Title, e.URL, e.Thumbnail, e.Views, e.LenVideosSugestions,
	)
}

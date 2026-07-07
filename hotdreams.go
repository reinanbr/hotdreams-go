// Package hotdreams is a Go port of the Python "hotdreams" library: an easy
// way to get video metadata (title, duration, thumbnail, preview, views,
// rating...) from several adult video sites.
//
// Each site has its own subpackage under sites/, exposing the same shape as
// the Python modules: SearchPorn(query, ...) returns a core.DataVideos
// (.Videos, .LenVideos, .LenPages, ...), and most also expose
// GetVideoEmbed(url) for fetching a single video's embed data and
// suggestions. For example:
//
//	import "github.com/reinanbr/hot_dreams/sites/pornhub"
//
//	result, err := pornhub.SearchPorn("lorena aquino", 2, nil)
//	fmt.Println(result.LenVideos, "videos found")
//	for _, video := range result.Videos {
//	    fmt.Println(video.Title, video.URL)
//	}
package hotdreams

import "github.com/reinanbr/hot_dreams/core"

// Version and Author mirror hotdreams.__version__ / hotdreams.__author__.
const (
	Version = "0.0.4"
	Author  = "ReinanBr <slimchatuba@gmail.com>"
)

// AvailableSites lists the site subpackages this lib ships, in the same
// order as Python's `available_sites`. Import the corresponding
// sites/<name> subpackage to use one.
var AvailableSites = []string{
	"XVideos",
	"PornHub",
	"tnaflix",
	"NoodleMagazine",
	"ukdevilz",
	"pornone",
	"Eporner",
	"FPO.XXX",
}

// SetUserAgent overrides the User-Agent header used by every site scraper.
func SetUserAgent(userAgent string) {
	core.SetUserAgent(userAgent)
}

// SetProxies routes every site scraper's requests through the given
// proxy/proxies, e.g. SetProxies("http://user:pass@host:port", "").
// Pass "" for a scheme you don't want to route through a proxy.
func SetProxies(httpProxy, httpsProxy string) error {
	return core.SetProxies(httpProxy, httpsProxy)
}

// ClearProxies stops routing requests through a proxy.
func ClearProxies() {
	core.ClearProxies()
}

// SetCookiesDir overrides the folder every site scraper stores/loads its
// cookie files in. Defaults to a folder inside the OS temp dir.
func SetCookiesDir(path string) error {
	return core.SetCookiesDir(path)
}

// GetDirectURL resolves a video page URL to the direct, playable/downloadable
// video file URL, using yt-dlp (must be installed and on PATH) under the
// hood. format defaults to "best" if empty.
func GetDirectURL(url string, format string) (string, error) {
	return core.GetDirectURL(url, format)
}

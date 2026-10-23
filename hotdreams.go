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
//	import "github.com/reinanbr/hotdreams-go/sites/pornhub"
//
//	result, err := pornhub.SearchPorn("lorena aquino", 2, nil)
//	fmt.Println(result.LenVideos, "videos found")
//	for _, video := range result.Videos {
//	    fmt.Println(video.Title, video.URL)
//	}
package hotdreams

import "github.com/reinanbr/hotdreams-go/core"

// Version and Author mirror hotdreams.__version__ / hotdreams.__author__.
const (
	Version = "0.3.0"
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
	"tubepornstars",
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

// SetProxy routes every request (http and https) through a single proxy,
// e.g. SetProxy("socks5://127.0.0.1:9050").
func SetProxy(proxy string) error {
	return core.SetProxy(proxy)
}

// SetProxyList rotates requests round-robin across the given proxies.
func SetProxyList(proxies ...string) error {
	return core.SetProxyList(proxies...)
}

// LoadProxyList reads proxies from a file (one per line) and rotates
// requests across them.
func LoadProxyList(path string) error {
	return core.LoadProxyList(path)
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

// ImportBrowserCookies loads cookies exported from a browser (Netscape
// cookies.txt or Cookie-Editor/EditThisCookie JSON) into every matching site
// session and saves them to that site's cookie file. Returns how many cookies
// were read.
func ImportBrowserCookies(path string) (int, error) {
	return core.ImportBrowserCookies(path)
}

// SetCookieHeader imports a raw "Cookie" header value copied from the browser
// devtools (e.g. "a=1; b=2") for the site at siteURL.
func SetCookieHeader(siteURL, header string) error {
	return core.SetCookieHeader(siteURL, header)
}

// GetDirectURL resolves a video page URL to the direct, playable/downloadable
// video file URL, using yt-dlp (must be installed and on PATH) under the
// hood. format defaults to "best" if empty.
func GetDirectURL(url string, format string) (string, error) {
	return core.GetDirectURL(url, format)
}

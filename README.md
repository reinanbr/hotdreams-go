<div align='center'>

<h1>Hot Dreams (Go)</h1>

<p>A Go library for scraping video metadata (title, duration, thumbnail, preview, views,
rating...) from several adult video sites — a Go port of
<a href="https://github.com/reinanbr/hotdreams">reinanbr/hotdreams</a> (Python).</p>

<img alt="License" src="https://img.shields.io/badge/license-BSD--3--Clause-orange">

</div>

<hr>

## Table of contents

- [Installation](#installation)
- [Quickstart](#quickstart)
- [Custom User-Agent and proxies](#custom-user-agent-and-proxies)
- [Custom cookies folder](#custom-cookies-folder)
- [Browser cookies](#browser-cookies)
- [Direct video URL](#direct-video-url)
- [Supported sites](#supported-sites)
- [License](#license)

<hr>

## Installation

```sh
go get github.com/reinanbr/hotdreams-go
```

Requires Go 1.21+.

## Quickstart

Every site has its own subpackage under `sites/`, each exposing the same shape as the
original Python modules: `SearchPorn(query, ...)` returns a `core.DataVideos`
(`.Videos`, `.LenVideos`, `.LenPages`, ...), and most also expose `GetVideoEmbed(url)` for
fetching a single video's embed data and suggestions.

```go
package main

import (
	"fmt"

	"github.com/reinanbr/hotdreams-go/sites/pornhub"
)

func main() {
	result, err := pornhub.SearchPorn("lorena aquino", 2, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.LenVideos, "videos found")

	for _, video := range result.Videos {
		fmt.Println(video.Title, video.URL)
	}
}
```

`SearchPorn(query, pageLimit, pageNumber)`: pass `pageLimit` pages to fetch (0 defaults to
2), or a non-nil `pageNumber` to fetch exactly that page instead.

<hr>

## Custom User-Agent and proxies

By default every site module ships with its own working headers. If you need to override
the User-Agent (e.g. to look like a different browser) or route requests through a proxy
(e.g. to avoid IP-based rate limiting), call these once — they apply to every site the
lib supports, retroactively:

```go
import "github.com/reinanbr/hotdreams-go"

hotdreams.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")

// one proxy for every request (http, https or socks5)
hotdreams.SetProxy("socks5://127.0.0.1:9050")

// or a different proxy per scheme
hotdreams.SetProxies("http://user:pass@host:port", "http://user:pass@host:port")

// or rotate round-robin across a list (one proxy per request)
hotdreams.SetProxyList("http://user:pass@host1:port", "socks5://host2:1080")
hotdreams.LoadProxyList("proxies.txt") // one proxy per line, '#' comments allowed

// ... use the lib normally
pornhub.SearchPorn("lorena aquino", 2, nil)

// to stop using the proxy
hotdreams.ClearProxies()
```

SOCKS5 proxies resolve hostnames on the proxy side (`socks5h://` is accepted too), which
also gets around sites blocked at the DNS level by your ISP.

<hr>

## Custom cookies folder

Every site module persists the cookies it picks up (e.g. anti-bot/session cookies) to a
small JSON file per site, and reuses them on the next call. By default those files are
written to a `hotdreams/cookies` folder inside the OS temp dir. Call `SetCookiesDir()` to
store them somewhere else instead:

```go
import "github.com/reinanbr/hotdreams-go"

hotdreams.SetCookiesDir("/home/me/.hotdreams_cookies")
```

<hr>

## Browser cookies

Some sites only serve full results to a "real" browser session (age gate, login,
Cloudflare challenge...). Export the cookies from your browser and import them; each
cookie goes to the site it belongs to and is saved to that site's cookie file, so it is
reused on the next runs too:

```go
import "github.com/reinanbr/hotdreams-go"

// Netscape cookies.txt (e.g. "Get cookies.txt LOCALLY" extension, yt-dlp --cookies)
// or the JSON export of Cookie-Editor / EditThisCookie — detected automatically
n, err := hotdreams.ImportBrowserCookies("/home/me/Downloads/cookies.txt")

// or paste the "Cookie" request header straight from the devtools Network tab
hotdreams.SetCookieHeader("https://pornone.com", "cf_clearance=...; age_verified=1")
```

If you import Cloudflare cookies (`cf_clearance`), also call `SetUserAgent` with the exact
User-Agent of the browser they came from, otherwise Cloudflare rejects them.

<hr>

## Direct video URL

`hotdreams.GetDirectURL(url, format)` resolves a video page URL to the direct,
playable/downloadable video file URL. It shells out to the
[yt-dlp](https://github.com/yt-dlp/yt-dlp) CLI, so make sure `yt-dlp` is installed and on
`PATH`:

```go
import "github.com/reinanbr/hotdreams-go"

direct, err := hotdreams.GetDirectURL("https://noodlemagazine.com/watch/-154109474_456276135", "best")
// direct -> "https://cdn77-pr.pvvstream.pro/videos/720/-154109474_456276135.mp4?..."
```

Supported hosts today: `pornhub.com`, `xvideos.com`, `tnaflix.com`, `spankbang.com`,
`playvids.com`, `noodlemagazine.com`, `eporner.com` — these are the sites yt-dlp ships a
dedicated extractor for. Calling it on any other host raises an error; yt-dlp's generic
extractor doesn't reliably resolve those, so they'd each need dedicated scraping, same as
the rest of this lib.

<hr>

## Supported sites

| site         | package                              | data video | embed url video |
|--------------|---------------------------------------|:----------:|:----------------:|
| PornHub       | `sites/pornhub`                       | ✅ | ✅ |
| XVideos       | `sites/xvideos`                       | ✅ | ✅ |
| tnaflix       | `sites/tnaflix`                       | ✅ | ✅ |
| NoodleMagazine| `sites/noodlemagazine`                | ✅ | ✅ |
| ukdevilz      | `sites/ukdevilz`                      | ✅ | ✅ |
| pornone       | `sites/pornone`                       | ✅ | ✅ |
| Eporner       | `sites/eporner`                       | ✅ | ✅ |
| FPO.XXX       | `sites/fpo`                           | ✅ | ✅ |
| SpankBang¹    | `sites/spankbang`                     | ✅ | ✅ |
| playvids²     | `sites/playvids`                      | ✅ |    |
| tubepornstars³| `sites/tubepornstars`                 | ✅ |    |

And sites from OnlyFans leaks:

| site         | package                              | data video |
|--------------|---------------------------------------|:----------:|
| nudes7        | `sites/only/nudes7`                   | ✅ |
| viralpornhub  | `sites/only/viralpornhub`             | ✅ |

¹ spankbang.com is behind Cloudflare's JS challenge; the upstream Python lib works around
it with `cloudscraper`, this port uses a plain HTTP client so requests may be blocked.

² playvids' `SearchPorn` returns a plain `[]Video` (no pagination metadata), matching the
upstream Python module's shape (it never went through the shared search helper there
either).

³ tubepornstars is an aggregator: each result's `URL` is the original video page on the
tube that hosts it (decoded from the site's `/out/` redirect link), and `Stats` holds that
tube's name. A pornstar name as the query lands on that pornstar's page.

<hr>

## License

BSD 3-Clause — see [LICENSE](LICENSE).

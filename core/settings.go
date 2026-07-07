package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Verbose gates the progress messages every site scraper prints while
// paginating (mirrors the `puts(...)` calls in the Python lib). Off by default.
var Verbose = false

func puts(format string, args ...any) {
	if Verbose {
		fmt.Printf(format+"\n", args...)
	}
}

// DefaultHeaders are the headers every site session ships with by default,
// copied verbatim from hotdreams.utils.settings.headers.
var DefaultHeaders = map[string]string{
	"Connection":         "keep-alive",
	"Cache-Control":      "max-age=0",
	"sec-ch-ua":          `"Chromium";v="110", "Not A(Brand";v="24", "Google Chrome";v="110"`,
	"sec-ch-ua-mobile":   "?0",
	"sec-ch-ua-platform": `"Linux"`,
	"User-Agent":         "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Accept":             "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
	"sec-fetch-site":     "none",
	"sec-fetch-mode":     "navigate",
	"sec-fetch-user":     "?1",
	"sec-fetch-dest":     "document",
	"Accept-Language":    "en-US,en;q=0.9,pt-BR;q=0.8,pt;q=0.7",
}

var (
	mu              sync.RWMutex
	cookiesDir      string
	customUserAgent string
	proxyHTTP       *url.URL
	proxyHTTPS      *url.URL
)

func init() {
	dir := filepath.Join(os.TempDir(), "hotdreams", "cookies")
	os.MkdirAll(dir, 0o755)
	cookiesDir = dir
}

// SetCookiesDir overrides the folder every site scraper stores/loads its
// cookie files in. Defaults to a folder inside the OS temp dir.
func SetCookiesDir(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return err
	}
	mu.Lock()
	cookiesDir = path
	mu.Unlock()
	return nil
}

// CookiesPath joins the current cookies dir with filename.
func CookiesPath(filename string) string {
	mu.RLock()
	defer mu.RUnlock()
	return filepath.Join(cookiesDir, filename)
}

// CookiesFileExists reports whether filename exists inside the cookies dir.
func CookiesFileExists(filename string) bool {
	_, err := os.Stat(CookiesPath(filename))
	return err == nil
}

// SetUserAgent overrides the User-Agent header used by every site scraper,
// retroactively (every Client reads this at request time).
func SetUserAgent(userAgent string) {
	mu.Lock()
	customUserAgent = userAgent
	mu.Unlock()
}

// SetProxies routes every site scraper's requests through the given
// proxy/proxies, e.g. SetProxies("http://user:pass@host:port", "").
func SetProxies(httpProxy, httpsProxy string) error {
	var hp, hsp *url.URL
	if httpProxy != "" {
		u, err := url.Parse(httpProxy)
		if err != nil {
			return err
		}
		hp = u
	}
	if httpsProxy != "" {
		u, err := url.Parse(httpsProxy)
		if err != nil {
			return err
		}
		hsp = u
	}
	mu.Lock()
	proxyHTTP = hp
	proxyHTTPS = hsp
	mu.Unlock()
	return nil
}

// ClearProxies stops routing requests through a proxy.
func ClearProxies() {
	mu.Lock()
	proxyHTTP = nil
	proxyHTTPS = nil
	mu.Unlock()
}

var sharedTransport = &http.Transport{
	Proxy: func(req *http.Request) (*url.URL, error) {
		mu.RLock()
		defer mu.RUnlock()
		if req.URL.Scheme == "https" && proxyHTTPS != nil {
			return proxyHTTPS, nil
		}
		if req.URL.Scheme == "http" && proxyHTTP != nil {
			return proxyHTTP, nil
		}
		if proxyHTTPS != nil {
			return proxyHTTPS, nil
		}
		return proxyHTTP, nil
	},
}

type headerRoundTripper struct{}

func (headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	mu.RLock()
	ua := customUserAgent
	mu.RUnlock()

	r2 := req.Clone(req.Context())
	for k, v := range DefaultHeaders {
		if r2.Header.Get(k) == "" {
			r2.Header.Set(k, v)
		}
	}
	if ua != "" {
		r2.Header.Set("User-Agent", ua)
	}
	return sharedTransport.RoundTrip(r2)
}

// Client is a per-site HTTP session: shared headers/proxy behavior (via the
// package-level settings above) plus its own cookie jar persisted to its own
// cookie file, mirroring one Python site module's `br`/`asession` + cookies file.
type Client struct {
	HTTP       *http.Client
	CookieFile string
	baseURL    string
}

// NewClient creates a session for a site, named cookieFile (e.g.
// ".cookies_PornHub.json"), and preloads any cookies already saved for it.
// baseURL is used as the cookie jar's reference origin for load/save.
func NewClient(cookieFile, baseURL string) *Client {
	jar, _ := cookiejar.New(nil)
	c := &Client{
		HTTP:       &http.Client{Transport: headerRoundTripper{}, Jar: jar},
		CookieFile: cookieFile,
		baseURL:    baseURL,
	}
	if CookiesFileExists(cookieFile) {
		_ = c.LoadCookies()
	}
	return c
}

// LoadCookies reads the client's cookie file (if any) and seeds the jar with it.
func (c *Client) LoadCookies() error {
	data, err := os.ReadFile(CookiesPath(c.CookieFile))
	if err != nil {
		return err
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return err
	}
	cookies := make([]*http.Cookie, 0, len(m))
	for k, v := range m {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	c.HTTP.Jar.SetCookies(u, cookies)
	return nil
}

// SaveCookies dumps the jar's cookies for finalURL to the client's cookie file.
func (c *Client) SaveCookies(finalURL string) error {
	u, err := url.Parse(finalURL)
	if err != nil {
		return err
	}
	cookies := c.HTTP.Jar.Cookies(u)
	m := make(map[string]string, len(cookies))
	for _, ck := range cookies {
		m[ck.Name] = ck.Value
	}
	data, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(CookiesPath(c.CookieFile), data, 0o644)
}

// Fetch performs a GET and returns the raw body, the final (post-redirect)
// URL, and the status code, saving cookies afterwards.
func (c *Client) Fetch(rawURL string) (body []byte, finalURL string, statusCode int, err error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", 0, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, "", 0, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", resp.StatusCode, err
	}
	finalURL = resp.Request.URL.String()
	_ = c.SaveCookies(finalURL)
	return body, finalURL, resp.StatusCode, nil
}

// pyQuote mimics Python's urllib.parse.quote(s) with the default safe="/":
// letters, digits, "_.-~" and "/" pass through unescaped, everything else
// (including spaces, which become %20) is percent-encoded.
func pyQuote(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') ||
			c == '_' || c == '.' || c == '-' || c == '~' || c == '/' {
			b.WriteByte(c)
		} else {
			fmt.Fprintf(&b, "%%%02X", c)
		}
	}
	return b.String()
}

// FetchPageFunc fetches and parses a single search-results page. It returns
// the parsed videos, or ok=false to signal "no more results" (mirrors the
// Python callback returning a falsy value to stop pagination).
type FetchPageFunc func(rawURL string, pageNumber int, query string) (videos []VideoData, ok bool, err error)

// BuildURLFunc builds the search URL for the given (already percent-encoded)
// query and page number.
type BuildURLFunc func(query string, page int) string

// SearchPornBase drives the same pagination loop as
// hotdreams.utils.settings.search_porn_base: fetch page 1..N (or just one
// page if pageNumber is given), stopping early the first time fetchPage
// signals no more results, and returns the aggregated DataVideos.
func SearchPornBase(query, urlBase, siteName string, pageLimit int, pageNumber *int, fetchPage FetchPageFunc, buildURL BuildURLFunc) (DataVideos, error) {
	quotedQuery := pyQuote(query)

	var videos []VideoData
	start := time.Now()
	var lastPage int

	if pageNumber == nil {
		n := pageLimit
		if n == 0 {
			n = 2
		}
		p := 1
		for i := 0; i < n; i++ {
			u := buildURL(quotedQuery, p)
			vds, ok, err := fetchPage(u, p, query)
			if err != nil {
				return DataVideos{}, err
			}
			if !ok || len(vds) == 0 {
				puts("\nops! finished in %d pages search!", p)
				break
			}
			videos = append(videos, vds...)
			puts("finding %d videos from %d pages of %d limit pages...", len(videos), p, n)
			p++
		}
		lastPage = p - 1
	} else {
		p := *pageNumber
		u := buildURL(quotedQuery, p)
		vds, _, err := fetchPage(u, p, query)
		if err != nil {
			return DataVideos{}, err
		}
		videos = vds
		puts("%d videos finds from page search number %d!", len(videos), p)
	}

	ping := time.Since(start).Seconds()
	lenVideos := len(videos)

	var lenPages *int
	var videosPerPage int
	if pageNumber != nil {
		videosPerPage = lenVideos
	} else {
		lp := lastPage
		lenPages = &lp
		if lp != 0 {
			videosPerPage = lenVideos / lp
		} else {
			videosPerPage = lenVideos
		}
	}

	return DataVideos{
		SiteName:       siteName,
		Success:        lenVideos > 0,
		Query:          quotedQuery,
		URLBase:        urlBase,
		URLSearch:      buildURL(quotedQuery, 1),
		Videos:         videos,
		Ping:           ping,
		VideosPerPages: videosPerPage,
		LenPages:       lenPages,
		LenVideos:      lenVideos,
	}, nil
}

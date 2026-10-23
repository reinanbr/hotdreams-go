package core

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	clientsMu sync.Mutex
	clients   []*Client
	// imported holds every browser cookie imported so far, so clients created
	// after the import still pick them up.
	imported []*http.Cookie
)

func registerClient(c *Client) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients = append(clients, c)
	c.applyCookies(imported)
}

// ImportBrowserCookies loads cookies exported from a browser and hands each
// one to the site session(s) whose domain it belongs to, persisting them to
// that site's cookie file. Useful to reuse a logged-in session or anti-bot
// clearance cookies (e.g. Cloudflare's cf_clearance — in that case also call
// SetUserAgent with the same User-Agent as the browser that got them).
//
// Supported formats (detected automatically):
//   - Netscape/Mozilla cookies.txt (yt-dlp, curl, "Get cookies.txt" extensions)
//   - JSON array exported by Cookie-Editor / EditThisCookie
//
// It returns how many cookies were imported.
func ImportBrowserCookies(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	var cookies []*http.Cookie
	if t := bytes.TrimSpace(data); len(t) > 0 && (t[0] == '[' || t[0] == '{') {
		cookies, err = ParseJSONCookies(t)
	} else {
		cookies, err = ParseNetscapeCookies(bytes.NewReader(data))
	}
	if err != nil {
		return 0, err
	}
	AddCookies(cookies)
	return len(cookies), nil
}

// SetCookieHeader imports a raw "Cookie" request header value (as copied from
// the browser devtools, e.g. "a=1; b=2") for the site at siteURL.
func SetCookieHeader(siteURL, header string) error {
	u, err := url.Parse(siteURL)
	if err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("invalid site url %q", siteURL)
	}
	parsed, err := http.ParseCookie(header)
	if err != nil {
		return err
	}
	for _, ck := range parsed {
		ck.Domain = u.Hostname()
		ck.Path = "/"
	}
	AddCookies(parsed)
	return nil
}

// AddCookies hands cookies (which must carry a Domain) to every site session
// whose host matches, and saves them to that site's cookie file.
func AddCookies(cookies []*http.Cookie) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	imported = append(imported, cookies...)
	for _, c := range clients {
		if c.applyCookies(cookies) {
			_ = c.SaveCookies(c.baseURL)
		}
	}
}

// applyCookies sets the cookies matching c's site into its jar, reporting
// whether any matched.
func (c *Client) applyCookies(cookies []*http.Cookie) bool {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return false
	}
	host := trimWWW(strings.ToLower(base.Hostname()))
	any := false
	for _, ck := range cookies {
		domain := trimWWW(strings.TrimPrefix(strings.ToLower(ck.Domain), "."))
		if domain == "" || !(domain == host || strings.HasSuffix(host, "."+domain) || strings.HasSuffix(domain, "."+host)) {
			continue
		}
		if !ck.Expires.IsZero() && ck.Expires.Before(time.Now()) {
			continue
		}
		// Set it against the client's own origin (host-only) so it is sent
		// regardless of www./non-www. and survives SaveCookies/LoadCookies.
		cp := *ck
		cp.Domain = ""
		if cp.Path == "" {
			cp.Path = "/"
		}
		c.HTTP.Jar.SetCookies(base, []*http.Cookie{&cp})
		any = true
	}
	return any
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}

// ParseNetscapeCookies parses a Netscape/Mozilla cookies.txt file:
// domain, include-subdomains, path, secure, expiry, name, value (tab separated).
func ParseNetscapeCookies(r io.Reader) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := strings.TrimRight(sc.Text(), "\r")
		httpOnly := false
		if strings.HasPrefix(line, "#HttpOnly_") {
			line = strings.TrimPrefix(line, "#HttpOnly_")
			httpOnly = true
		}
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}
		f := strings.Split(line, "\t")
		if len(f) < 7 {
			continue
		}
		ck := &http.Cookie{
			Domain:   f[0],
			Path:     f[2],
			Secure:   strings.EqualFold(f[3], "TRUE"),
			Name:     f[5],
			Value:    f[6],
			HttpOnly: httpOnly,
		}
		if exp, err := strconv.ParseInt(f[4], 10, 64); err == nil && exp > 0 {
			ck.Expires = time.Unix(exp, 0)
		}
		cookies = append(cookies, ck)
	}
	return cookies, sc.Err()
}

// ParseJSONCookies parses the JSON export of the Cookie-Editor /
// EditThisCookie browser extensions (an array of cookie objects).
func ParseJSONCookies(data []byte) ([]*http.Cookie, error) {
	var raw []struct {
		Name           string  `json:"name"`
		Value          string  `json:"value"`
		Domain         string  `json:"domain"`
		Path           string  `json:"path"`
		Secure         bool    `json:"secure"`
		HTTPOnly       bool    `json:"httpOnly"`
		ExpirationDate float64 `json:"expirationDate"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("unsupported cookies JSON (expected an array of cookies): %w", err)
	}
	cookies := make([]*http.Cookie, 0, len(raw))
	for _, r := range raw {
		ck := &http.Cookie{
			Name:     r.Name,
			Value:    r.Value,
			Domain:   r.Domain,
			Path:     r.Path,
			Secure:   r.Secure,
			HttpOnly: r.HTTPOnly,
		}
		if r.ExpirationDate > 0 {
			ck.Expires = time.Unix(int64(r.ExpirationDate), 0)
		}
		cookies = append(cookies, ck)
	}
	return cookies, nil
}

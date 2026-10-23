package core

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func cookieValue(c *Client, rawURL, name string) string {
	u, _ := url.Parse(rawURL)
	for _, ck := range c.HTTP.Jar.Cookies(u) {
		if ck.Name == name {
			return ck.Value
		}
	}
	return ""
}

func TestImportBrowserCookiesNetscape(t *testing.T) {
	SetCookiesDir(t.TempDir())
	c := NewClient(".cookies_test_ns.json", "https://www.example-site.com")
	other := NewClient(".cookies_test_other.json", "https://other.com")

	path := filepath.Join(t.TempDir(), "cookies.txt")
	content := "# Netscape HTTP Cookie File\n" +
		".example-site.com\tTRUE\t/\tTRUE\t0\tsession\tabc\n" +
		"#HttpOnly_example-site.com\tFALSE\t/\tFALSE\t9999999999\tcf_clearance\txyz\n" +
		".example-site.com\tTRUE\t/\tFALSE\t1\texpired\told\n" +
		".unrelated.com\tTRUE\t/\tFALSE\t0\tnope\t1\n"
	os.WriteFile(path, []byte(content), 0o644)

	n, err := ImportBrowserCookies(path)
	if err != nil || n != 4 {
		t.Fatalf("ImportBrowserCookies = %d, %v", n, err)
	}
	if got := cookieValue(c, "https://www.example-site.com/x", "session"); got != "abc" {
		t.Errorf("session = %q", got)
	}
	if got := cookieValue(c, "https://www.example-site.com/", "cf_clearance"); got != "xyz" {
		t.Errorf("cf_clearance = %q", got)
	}
	if got := cookieValue(c, "https://www.example-site.com/", "expired"); got != "" {
		t.Errorf("expired cookie was imported: %q", got)
	}
	if got := cookieValue(other, "https://other.com/", "session"); got != "" {
		t.Errorf("cookie leaked to other site: %q", got)
	}
	data, _ := os.ReadFile(CookiesPath(".cookies_test_ns.json"))
	if !strings.Contains(string(data), "cf_clearance") {
		t.Errorf("cookie file not persisted: %s", data)
	}
}

func TestImportBrowserCookiesJSON(t *testing.T) {
	SetCookiesDir(t.TempDir())
	c := NewClient(".cookies_test_json.json", "https://jsonsite.com")
	path := filepath.Join(t.TempDir(), "cookies.json")
	os.WriteFile(path, []byte(`[{"name":"tok","value":"42","domain":".jsonsite.com","path":"/","expirationDate":9999999999.5}]`), 0o644)
	if _, err := ImportBrowserCookies(path); err != nil {
		t.Fatal(err)
	}
	if got := cookieValue(c, "https://jsonsite.com/", "tok"); got != "42" {
		t.Errorf("tok = %q", got)
	}
}

func TestSetCookieHeader(t *testing.T) {
	SetCookiesDir(t.TempDir())
	c := NewClient(".cookies_test_hdr.json", "https://hdrsite.com")
	if err := SetCookieHeader("https://hdrsite.com", "a=1; b=2"); err != nil {
		t.Fatal(err)
	}
	if cookieValue(c, "https://hdrsite.com/", "a") != "1" || cookieValue(c, "https://hdrsite.com/", "b") != "2" {
		t.Error("header cookies not applied")
	}
}

func TestProxySelection(t *testing.T) {
	defer ClearProxies()
	req, _ := http.NewRequest("GET", "https://x.com", nil)

	if err := SetProxy("socks5h://127.0.0.1:9050"); err != nil {
		t.Fatal(err)
	}
	if u, _ := selectProxy(req); u.String() != "socks5://127.0.0.1:9050" {
		t.Errorf("SetProxy -> %v", u)
	}
	if err := SetProxy("ftp://bad:1"); err == nil {
		t.Error("expected error for ftp scheme")
	}

	SetProxyList("http://a:1", "http://b:2")
	var got []string
	for i := 0; i < 3; i++ {
		u, _ := selectProxy(req)
		got = append(got, u.Host)
	}
	if strings.Join(got, ",") != "a:1,b:2,a:1" {
		t.Errorf("rotation = %v", got)
	}

	ClearProxies()
	if u, _ := selectProxy(req); u != nil {
		t.Errorf("expected no proxy, got %v", u)
	}
}

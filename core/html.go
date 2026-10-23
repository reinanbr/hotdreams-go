package core

import (
	"bytes"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ParseHTML parses a raw HTML body into a goquery document, the Go
// equivalent of `bs(html, features="html.parser")`.
func ParseHTML(body []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewReader(body))
}

// Meta returns the content attribute of the first <meta property="prop">
// (or name="prop") tag found, "" if absent.
func Meta(doc *goquery.Document, key, prop string) string {
	sel := doc.Find("meta[" + key + "=\"" + prop + "\"]").First()
	if sel.Length() == 0 {
		return ""
	}
	v, _ := sel.Attr("content")
	return v
}

// Attr returns attr of the first match of selector, "" if absent.
func Attr(doc *goquery.Document, selector, attr string) string {
	sel := doc.Find(selector).First()
	if sel.Length() == 0 {
		return ""
	}
	v, _ := sel.Attr(attr)
	return v
}

// Text returns the trimmed text of the first match of selector, "" if absent.
func Text(doc *goquery.Document, selector string) string {
	sel := doc.Find(selector).First()
	if sel.Length() == 0 {
		return ""
	}
	return strings.TrimSpace(sel.Text())
}

// lazyImageAttrs are checked, in order, before falling back to src: most
// sites lazy-load thumbnails, leaving a placeholder (data: URI, blank gif)
// in src and the real image URL in one of these attributes.
var lazyImageAttrs = []string{"data-src", "data-original", "data-lazy-src", "data-thumb", "data-srcset", "srcset", "src"}

// ImageURL returns the real image URL of an <img> (or <source>) selection,
// skipping lazy-load placeholders and resolving protocol-relative/relative
// URLs against base. Returns "" if no usable URL is found.
func ImageURL(img *goquery.Selection, base string) string {
	for _, attr := range lazyImageAttrs {
		v, ok := img.Attr(attr)
		if !ok {
			continue
		}
		v = strings.TrimSpace(v)
		if strings.HasSuffix(attr, "srcset") {
			// "url1 320w, url2 640w" -> url1
			v = strings.TrimSpace(strings.Split(v, ",")[0])
			if f := strings.Fields(v); len(f) > 0 {
				v = f[0]
			}
		}
		if v == "" || strings.HasPrefix(v, "data:") || isPlaceholderImage(v) {
			continue
		}
		return ResolveURL(base, v)
	}
	return ""
}

func isPlaceholderImage(v string) bool {
	l := strings.ToLower(v)
	return strings.Contains(l, "blank.gif") || strings.Contains(l, "placeholder") ||
		strings.Contains(l, "lazy.") || strings.HasSuffix(l, "/1x1.gif")
}

// ResolveURL resolves ref (relative or protocol-relative) against base.
func ResolveURL(base, ref string) string {
	if strings.HasPrefix(ref, "//") {
		return "https:" + ref
	}
	b, err := url.Parse(base)
	if err != nil || base == "" {
		return ref
	}
	r, err := url.Parse(ref)
	if err != nil {
		return ref
	}
	return b.ResolveReference(r).String()
}

// ParseDurationClock converts a "H:MM:SS" / "MM:SS" / "SS" clock string
// (used across most sites for the video duration label) into seconds.
func ParseDurationClock(clock string) (int, bool) {
	clock = strings.TrimSpace(clock)
	if clock == "" {
		return 0, false
	}
	parts := strings.Split(clock, ":")
	nums := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return 0, false
		}
		nums = append(nums, n)
	}
	switch len(nums) {
	case 3:
		return nums[0]*3600 + nums[1]*60 + nums[2], true
	case 2:
		return nums[0]*60 + nums[1], true
	case 1:
		return nums[0], true
	default:
		return 0, false
	}
}

// AtoiOr parses s as an int, returning def on failure.
func AtoiOr(s string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return n
}

package core

import (
	"bytes"
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

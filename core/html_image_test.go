package core

import "testing"

func TestImageURLLazy(t *testing.T) {
	doc, _ := ParseHTML([]byte(`
<img id="a" src="data:image/gif;base64,R0lGOD" data-src="//cdn.x.com/a.jpg">
<img id="b" src="/img/blank.gif" data-original="https://cdn.x.com/b.jpg">
<img id="c" src="/thumbs/c.jpg">
<img id="d" srcset="https://cdn.x.com/d1.jpg 320w, https://cdn.x.com/d2.jpg 640w">`))
	cases := map[string]string{
		"#a": "https://cdn.x.com/a.jpg",
		"#b": "https://cdn.x.com/b.jpg",
		"#c": "https://site.com/thumbs/c.jpg",
		"#d": "https://cdn.x.com/d1.jpg",
	}
	for sel, want := range cases {
		if got := ImageURL(doc.Find(sel), "https://site.com"); got != want {
			t.Errorf("%s: got %q, want %q", sel, got, want)
		}
	}
}

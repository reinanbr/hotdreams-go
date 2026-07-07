package core

import "testing"

func TestParseDurationClock(t *testing.T) {
	cases := []struct {
		in     string
		want   int
		wantOK bool
	}{
		{"1:02:03", 3723, true},
		{"12:37", 757, true},
		{"45", 45, true},
		{"", 0, false},
		{"not-a-clock", 0, false},
	}
	for _, c := range cases {
		got, ok := ParseDurationClock(c.in)
		if ok != c.wantOK || (ok && got != c.want) {
			t.Errorf("ParseDurationClock(%q) = (%d, %v), want (%d, %v)", c.in, got, ok, c.want, c.wantOK)
		}
	}
}

func TestAtoiOr(t *testing.T) {
	if got := AtoiOr("42", -1); got != 42 {
		t.Errorf("AtoiOr(42) = %d, want 42", got)
	}
	if got := AtoiOr("nope", -1); got != -1 {
		t.Errorf("AtoiOr(nope) = %d, want -1", got)
	}
}

package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strings"
)

// YtdlpSupportedHosts are the hosts yt-dlp ships a dedicated extractor for,
// verified against this lib's sites (mirrors utils.direct_url.YTDLP_SUPPORTED_HOSTS).
var YtdlpSupportedHosts = []string{
	"pornhub.com",
	"xvideos.com",
	"tnaflix.com",
	"spankbang.com",
	"playvids.com",
	"noodlemagazine.com",
	"eporner.com",
}

func hostOf(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := strings.ToLower(u.Hostname())
	return strings.TrimPrefix(host, "www.")
}

func isSupportedHost(host string) bool {
	for _, h := range YtdlpSupportedHosts {
		if host == h || strings.HasSuffix(host, "."+h) {
			return true
		}
	}
	return false
}

// GetDirectURL resolves a video page URL to the direct, playable/downloadable
// video file URL, for any host in YtdlpSupportedHosts. It shells out to the
// yt-dlp binary (the Go equivalent of the Python lib's optional `yt-dlp`
// dependency) — install it separately (e.g. `pip install yt-dlp` or your
// package manager) for this to work.
func GetDirectURL(rawURL string, format string) (string, error) {
	if format == "" {
		format = "best"
	}
	host := hostOf(rawURL)
	if !isSupportedHost(host) {
		return "", fmt.Errorf(
			"no direct-url extractor for host '%s' yet. Supported hosts: %s",
			host, strings.Join(YtdlpSupportedHosts, ", "),
		)
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp is required for GetDirectURL(); install it and make sure it's on PATH: %w", err)
	}

	cmd := exec.Command("yt-dlp", "--quiet", "--no-warnings", "-f", format, "--skip-download", "--dump-json", rawURL)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %w: %s", err, stderr.String())
	}

	var info struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &info); err != nil {
		return "", fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}
	if info.URL == "" {
		return "", fmt.Errorf("yt-dlp returned no direct url for %s", rawURL)
	}
	return info.URL, nil
}

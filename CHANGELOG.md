# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-10-23

### Added
- `sites/tubepornstars`: search for tubepornstars.com. Each result's `URL` is the original
  video page on the hosting tube, decoded from the site's `/out/` redirect link; pornstar
  name queries follow the redirect to the pornstar page, with pagination.
- Proxy helpers: `SetProxy` (single proxy for http and https), `SetProxyList` /
  `LoadProxyList` (round-robin rotation), and SOCKS5 support (`socks5://`, `socks5h://`).
- Browser cookies: `ImportBrowserCookies` (Netscape cookies.txt or Cookie-Editor /
  EditThisCookie JSON) and `SetCookieHeader` (raw `Cookie` header from devtools).
- `core.ImageURL` / `core.ResolveURL` helpers for lazy-loaded thumbnails and relative URLs.

### Fixed
- pornone: thumbnails from lazy-loaded images and relative video URLs.

## [0.2.0]

### Changed
- Module renamed to `github.com/reinanbr/hotdreams-go`.

## [0.1.0]

### Added
- First release: Go port of the Python hotdreams library, with CI for tests and tagged
  releases.

[0.3.0]: https://github.com/reinanbr/hotdreams-go/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/reinanbr/hotdreams-go/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/reinanbr/hotdreams-go/releases/tag/v0.1.0

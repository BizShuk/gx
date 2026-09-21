package applepodcast

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// Target 是解析使用者輸入後得到的查詢對象：節目的 collection ID。
type Target struct {
	ID string
}

// idInPath 比對路徑上的 `id<數字>` 段，如 /tw/podcast/<slug>/id1702409419。
var idInPath = regexp.MustCompile(`(?:^|/)id(\d{1,18})(?:/|$)`)

// ParseTarget 解析節目網址，只接受 Apple 網域的網址：
//
//	https://podcasts.apple.com/tw/podcast/<slug>/id1702409419?l=en-GB
//	https://itunes.apple.com/us/podcast/id1702409419
//	podcasts.apple.com/us/podcast/<slug>/id1702409419（省略 scheme）
//
// 裸 collection ID 不收：一串數字說不出它屬於哪個平台，
// 網址才是使用者手上真正有的東西。
// 單集網址（?i=<episode>）仍指向節目本身，取路徑上的節目 ID。
func ParseTarget(input string) (Target, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return Target{}, fmt.Errorf("empty podcast target")
	}

	raw := s
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || !IsHost(u.Host) {
		return Target{}, fmt.Errorf("not an Apple Podcasts url: %q", input)
	}
	if m := idInPath.FindStringSubmatch(u.Path); len(m) == 2 {
		return Target{ID: m[1]}, nil
	}
	if q := u.Query().Get("id"); IsCollectionID(q) {
		return Target{ID: q}, nil
	}
	return Target{}, fmt.Errorf("no podcast id in %q", input)
}

// IsCollectionID 判斷輸入本身是否已是一個 collection ID（純數字、不以 0 開頭）。
func IsCollectionID(input string) bool {
	s := strings.TrimSpace(input)
	if s == "" || s[0] == '0' || len(s) > 18 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// IsHost 判斷 host 是否為 Apple Podcasts／iTunes 網域。
func IsHost(host string) bool {
	h := strings.ToLower(host)
	return h == "podcasts.apple.com" || h == "itunes.apple.com"
}

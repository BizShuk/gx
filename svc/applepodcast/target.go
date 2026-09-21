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

// ParseTarget 解析使用者輸入，接受下列寫法：
//
//	1702409419
//	id1702409419
//	https://podcasts.apple.com/tw/podcast/<slug>/id1702409419?l=en-GB
//	https://itunes.apple.com/us/podcast/id1702409419
//
// 單集網址（?i=<episode>）仍指向節目本身，取路徑上的節目 ID。
func ParseTarget(input string) (Target, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return Target{}, fmt.Errorf("empty podcast target")
	}

	if IsCollectionID(s) {
		return Target{ID: s}, nil
	}
	if len(s) > 2 && strings.EqualFold(s[:2], "id") && IsCollectionID(s[2:]) {
		return Target{ID: s[2:]}, nil
	}

	raw := s
	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil || !isAppleHost(u.Host) {
		return Target{}, fmt.Errorf("no podcast id in %q", input)
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

func isAppleHost(host string) bool {
	h := strings.ToLower(host)
	return h == "podcasts.apple.com" || h == "itunes.apple.com"
}

package bilibili

import (
	"fmt"
	"strings"
)

// Target 是解析使用者輸入後得到的查詢對象：已帶空間 UID 時填 ID。
type Target struct {
	ID string
}

// ParseTarget 解析使用者輸入，接受下列寫法：
//
//	2267573
//	https://space.bilibili.com/2267573
//	https://space.bilibili.com/2267573/video
//	https://m.bilibili.com/space/2267573
func ParseTarget(input string) (Target, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return Target{}, fmt.Errorf("empty space target")
	}

	if IsSpaceID(s) {
		return Target{ID: s}, nil
	}

	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	if i := strings.IndexAny(s, "?#"); i >= 0 {
		s = s[:i]
	}

	segments := splitPath(s)
	if len(segments) > 0 && IsHost(segments[0]) {
		host := strings.ToLower(segments[0])
		segments = segments[1:]
		if host == "b23.tv" || host == "b23.wtf" {
			return Target{}, fmt.Errorf("b23 short link is not a space id: %q", input)
		}
	}
	if len(segments) == 0 {
		return Target{}, fmt.Errorf("no space id in %q", input)
	}

	if IsSpaceID(segments[0]) {
		return Target{ID: segments[0]}, nil
	}

	for i, segment := range segments {
		if segment == "space" && i+1 < len(segments) && IsSpaceID(segments[i+1]) {
			return Target{ID: segments[i+1]}, nil
		}
	}

	return Target{}, fmt.Errorf("no space id in %q", input)
}

// IsSpaceID 判斷輸入本身是否已是一個空間 UID（純數字）。
func IsSpaceID(input string) bool {
	s := strings.TrimSpace(input)
	if s == "" || s[0] == '0' || len(s) > 16 {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func splitPath(s string) []string {
	var segments []string
	for segment := range strings.SplitSeq(s, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}
	return segments
}

// IsHost 判斷 host 是否為 Bilibili 網域（含 b23 短鏈網域）。
func IsHost(segment string) bool {
	host := strings.ToLower(segment)
	return host == "b23.tv" ||
		host == "b23.wtf" ||
		host == "bilibili.com" ||
		strings.HasSuffix(host, ".bilibili.com")
}

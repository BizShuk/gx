package youtube

import (
	"fmt"
	"strings"
)

// Target 是解析使用者輸入後得到的查詢對象：
// 已經帶頻道 ID 時填 ID，否則填 @handle。
type Target struct {
	ID     string
	Handle string
}

// ParseTarget 解析使用者輸入，接受下列寫法：
//
//	UCBR8-60-B28hp2BmDPdntcQ
//	@YouTube / YouTube
//	https://www.youtube.com/@YouTube
//	https://www.youtube.com/@YouTube/videos
//	https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ
//	youtube.com/c/YouTube、youtube.com/user/YouTube
func ParseTarget(input string) (Target, error) {
	s := strings.TrimSpace(input)
	if s == "" {
		return Target{}, fmt.Errorf("empty channel target")
	}

	if IsChannelID(s) {
		return Target{ID: s}, nil
	}

	// 去掉 scheme 與查詢字串／錨點，剩下 host + path 便於逐段解析。
	if i := strings.Index(s, "://"); i >= 0 {
		s = s[i+3:]
	}
	if i := strings.IndexAny(s, "?#"); i >= 0 {
		s = s[:i]
	}

	segments := splitPath(s)
	if len(segments) > 0 && IsHost(segments[0]) {
		if strings.EqualFold(segments[0], "youtu.be") {
			return Target{}, fmt.Errorf("video short link is not a channel: %q", input)
		}
		segments = segments[1:]
	}
	if len(segments) == 0 {
		return Target{}, fmt.Errorf("no channel handle in %q", input)
	}
	// 影片與播放清單的路徑不是頻道；不擋的話 `/playlist` 會被當成 handle `@playlist`。
	if nonChannelPaths[strings.ToLower(segments[0])] {
		return Target{}, fmt.Errorf("not a channel url: %q", input)
	}

	for i, segment := range segments {
		switch {
		case strings.HasPrefix(segment, "@"):
			return newHandleTarget(segment, input)
		case segment == "channel" && i+1 < len(segments):
			if id := segments[i+1]; IsChannelID(id) {
				return Target{ID: id}, nil
			}
			return Target{}, fmt.Errorf("invalid channel id in %q", input)
		case (segment == "c" || segment == "user") && i+1 < len(segments):
			return newHandleTarget(segments[i+1], input)
		}
	}

	// 沒有任何路徑關鍵字時，把第一段當成 handle（例如裸寫 `YouTube`）。
	return newHandleTarget(segments[0], input)
}

// IsChannelID 判斷輸入本身是否已是一個頻道 ID。
func IsChannelID(input string) bool {
	return channelIDPattern.MatchString(strings.TrimSpace(input))
}

func newHandleTarget(segment, input string) (Target, error) {
	handle := strings.TrimPrefix(segment, "@")
	if handle == "" {
		return Target{}, fmt.Errorf("no channel handle in %q", input)
	}
	return Target{Handle: "@" + handle}, nil
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

// nonChannelPaths 是 YouTube 網址中指向影片或播放清單、而非頻道的第一段路徑。
var nonChannelPaths = map[string]bool{
	"watch": true, "shorts": true, "live": true, "embed": true, "playlist": true,
}

// IsHost 判斷 host 是否為 YouTube 網域。
func IsHost(segment string) bool {
	host := strings.ToLower(segment)
	return host == "youtu.be" ||
		host == "youtube.com" ||
		strings.HasSuffix(host, ".youtube.com")
}

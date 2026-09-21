// Package resolve 是跨平台的頻道網址解析入口：依網址網域分派給對應平台。
//
// 它是 gx 對外的函式庫介面 —— 呼叫端拿到一個使用者貼上的網址，
// 不必先知道它屬於哪個平台。平台規則仍各自住在 svc/<platform>，這裡只做分派。
package resolve

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bizshuk/gx/model"
	"github.com/bizshuk/gx/svc/applepodcast"
	"github.com/bizshuk/gx/svc/bilibili"
	"github.com/bizshuk/gx/svc/youtube"
)

// ErrUnknownURL 代表輸入不是網址，或網域不屬於任何已支援的平台。
var ErrUnknownURL = errors.New("not a supported channel url")

// hosts 依序列出各平台的網域判斷；第一個吻合的平台勝出。
var hosts = []struct {
	platform string
	match    func(host string) bool
}{
	{model.PLATFORM_YOUTUBE, youtube.IsHost},
	{model.PLATFORM_BILIBILI, bilibili.IsHost},
	{model.PLATFORM_APPLE_PODCAST, applepodcast.IsHost},
}

// Platform 依網址網域判斷所屬平台（model.PLATFORM_*）。
// 不是網址或網域不認得時回空字串 —— 只看網域、不看 ID 形狀：
// 一串裸 ID 說不出它屬於哪個平台。
func Platform(input string) string {
	host := hostOf(input)
	if host == "" {
		return ""
	}
	for _, h := range hosts {
		if h.match(host) {
			return h.platform
		}
	}
	return ""
}

// Resolver 持有各平台的 client，依網址分派查詢。
type Resolver struct {
	youtube      *youtube.Client
	bilibili     *bilibili.Client
	applePodcast *applepodcast.Client
}

// Option 覆寫某個平台的 client，供函式庫呼叫端帶入自己的逾時與來源網域。
type Option func(*Resolver)

// WithYouTube 使用指定的 YouTube client。
func WithYouTube(c *youtube.Client) Option { return func(r *Resolver) { r.youtube = c } }

// WithBilibili 使用指定的 Bilibili client。
func WithBilibili(c *bilibili.Client) Option { return func(r *Resolver) { r.bilibili = c } }

// WithApplePodcast 使用指定的 Apple Podcasts client。
func WithApplePodcast(c *applepodcast.Client) Option {
	return func(r *Resolver) { r.applePodcast = c }
}

// New 建立 Resolver；未指定的平台使用該平台的預設 client。
func New(opts ...Option) *Resolver {
	r := &Resolver{}
	for _, opt := range opts {
		opt(r)
	}
	if r.youtube == nil {
		r.youtube = youtube.NewClient()
	}
	if r.bilibili == nil {
		r.bilibili = bilibili.NewClient()
	}
	if r.applePodcast == nil {
		r.applePodcast = applepodcast.NewClient()
	}
	return r
}

// Channel 依網址網域分派給對應平台，回傳標準 Channel。
func (r *Resolver) Channel(ctx context.Context, input string) (*model.Channel, error) {
	switch Platform(input) {
	case model.PLATFORM_YOUTUBE:
		return r.youtube.GetChannel(ctx, input)
	case model.PLATFORM_BILIBILI:
		return r.bilibili.GetChannel(ctx, input)
	case model.PLATFORM_APPLE_PODCAST:
		return r.applePodcast.GetChannel(ctx, input)
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownURL, input)
	}
}

// hostOf 取出網址的 host；省略 scheme 的 `youtube.com/@x` 也算網址。
// 沒有 host 或 host 不含 `.` 時回空字串（`@handle`、裸 ID 都落在這裡）。
func hostOf(input string) string {
	s := strings.TrimSpace(input)
	if s == "" {
		return ""
	}
	if !strings.Contains(s, "://") {
		s = "https://" + s
	}
	u, err := url.Parse(s)
	if err != nil || u.User != nil {
		return ""
	}
	host := strings.ToLower(u.Hostname())
	if !strings.Contains(host, ".") {
		return ""
	}
	return host
}

package applepodcast

import (
	"time"

	"github.com/spf13/viper"
)

// 設定 key：與 config/default_settings.json 的扁平欄位對應，
// 由本套件自行讀取，呼叫端不需要把設定值一路搬進來。
const (
	KEY_BASE_URL   = "apple_podcast_base_url"
	KEY_TIMEOUT    = "http_timeout"
	KEY_USER_AGENT = "http_user_agent"
)

type options struct {
	baseURL   string
	timeout   time.Duration
	userAgent string
}

// Option 以函式選項覆寫設定值，主要供測試指向本地伺服器。
type Option func(*options)

// WithBaseURL 覆寫 iTunes lookup 的來源網域。
func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		if baseURL != "" {
			o.baseURL = baseURL
		}
	}
}

// WithTimeout 覆寫單次 HTTP 請求的上限。
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) {
		if timeout > 0 {
			o.timeout = timeout
		}
	}
}

// WithUserAgent 覆寫請求的 User-Agent。
func WithUserAgent(userAgent string) Option {
	return func(o *options) {
		if userAgent != "" {
			o.userAgent = userAgent
		}
	}
}

// applyOptions 取值順序為：函式選項 > viper 設定 > 套件預設值。
func applyOptions(opts ...Option) *options {
	o := &options{
		baseURL:   DEFAULT_BASE_URL,
		timeout:   DEFAULT_TIMEOUT,
		userAgent: DEFAULT_USER_AGENT,
	}

	WithBaseURL(viper.GetString(KEY_BASE_URL))(o)
	WithTimeout(viper.GetDuration(KEY_TIMEOUT))(o)
	WithUserAgent(viper.GetString(KEY_USER_AGENT))(o)

	for _, opt := range opts {
		opt(o)
	}
	return o
}

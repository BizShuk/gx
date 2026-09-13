package bilibili

import "github.com/spf13/viper"

// 設定 key：與 config/default_settings.json 的扁平欄位對應，
// 由本套件自行讀取，呼叫端不需要把設定值一路搬進來。
const KEY_BASE_URL = "bilibili_base_url"

type options struct {
	baseURL string
}

// Option 以函式選項覆寫設定值，主要供測試指向本地基底網址。
type Option func(*options)

// WithBaseURL 覆寫 Bilibili 空間頁的來源網域。
func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		if baseURL != "" {
			o.baseURL = baseURL
		}
	}
}

// applyOptions 取值順序為：函式選項 > viper 設定 > 套件預設值。
func applyOptions(opts ...Option) *options {
	o := &options{baseURL: DEFAULT_BASE_URL}
	WithBaseURL(viper.GetString(KEY_BASE_URL))(o)
	for _, opt := range opts {
		opt(o)
	}
	return o
}

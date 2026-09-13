// Package bilibili 提供不需 API key 的 Bilibili 空間識別解析。
package bilibili

import "strings"

const (
	// DEFAULT_BASE_URL 是 Bilibili UP 主空間的來源網域。
	DEFAULT_BASE_URL = "https://space.bilibili.com"
)

// Client 組出空間的正規網址。UID 已在輸入裡時不發出請求。
type Client struct {
	baseURL string
}

// NewClient 建立 Client，未指定的選項套用套件預設值。
func NewClient(opts ...Option) *Client {
	o := applyOptions(opts...)
	return &Client{baseURL: strings.TrimRight(o.baseURL, "/")}
}

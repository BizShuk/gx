package bilibili

import (
	"context"
	"fmt"
)

// Channel 是一個 Bilibili UP 主空間的識別資訊。
// 平台沒有官方 channel RSS，因此沒有 rss 欄位——不代填第三方合成源。
type Channel struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// GetChannel 由 UID 或空間網址解析出正規的空間網址。
// UID 已在輸入裡時直接組出結果，不發出請求。
func (c *Client) GetChannel(_ context.Context, input string) (*Channel, error) {
	target, err := ParseTarget(input)
	if err != nil {
		return nil, err
	}
	if target.ID == "" {
		return nil, fmt.Errorf("no space id in %q", input)
	}
	return &Channel{ID: target.ID, URL: c.channelURL(target.ID)}, nil
}

func (c *Client) channelURL(id string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, id)
}

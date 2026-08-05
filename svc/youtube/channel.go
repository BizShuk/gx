package youtube

import (
	"context"
	"errors"
	"fmt"
)

// Channel 是一個 YouTube 頻道的識別資訊。
type Channel struct {
	Handle string `json:"handle,omitempty"`
	ID     string `json:"id"`
	URL    string `json:"url"`
}

// GetChannel 由 handle 或網址解析出頻道 ID 與正規網址。
// 輸入本身已帶頻道 ID 時直接組出結果，不發出請求。
func (c *Client) GetChannel(ctx context.Context, input string) (*Channel, error) {
	target, err := ParseTarget(input)
	if err != nil {
		return nil, err
	}
	if target.ID != "" {
		return &Channel{ID: target.ID, URL: c.channelURL(target.ID)}, nil
	}

	html, err := c.fetch(ctx, fmt.Sprintf("%s/%s", c.baseURL, target.Handle))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("channel %s not found", target.Handle)
		}
		return nil, err
	}

	id, ok := ExtractChannelID(html)
	if !ok {
		return nil, fmt.Errorf("no channel id found for %s", target.Handle)
	}

	return &Channel{Handle: target.Handle, ID: id, URL: c.channelURL(id)}, nil
}

func (c *Client) channelURL(id string) string {
	return fmt.Sprintf("%s/channel/%s", c.baseURL, id)
}

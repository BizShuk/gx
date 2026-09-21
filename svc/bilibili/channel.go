package bilibili

import (
	"context"
	"fmt"

	"github.com/bizshuk/gx/model"
)

// GetChannel 由 UID 或空間網址解析出正規的空間網址。
// UID 已在輸入裡時直接組出結果，不發出請求。
// 平台沒有官方 channel RSS，因此 RSS 留空 —— 不代填第三方合成源。
func (c *Client) GetChannel(_ context.Context, input string) (*model.Channel, error) {
	target, err := ParseTarget(input)
	if err != nil {
		return nil, err
	}
	if target.ID == "" {
		return nil, fmt.Errorf("no space id in %q", input)
	}
	return &model.Channel{
		Platform: model.PLATFORM_BILIBILI,
		ID:       target.ID,
		URL:      c.channelURL(target.ID),
	}, nil
}

func (c *Client) channelURL(id string) string {
	return fmt.Sprintf("%s/%s", c.baseURL, id)
}

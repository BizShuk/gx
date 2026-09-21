package youtube

import (
	"context"
	"errors"
	"fmt"

	"github.com/bizshuk/gx/model"
)

// GetChannel 由 handle、網址或頻道 ID 解析出頻道 ID、名稱、正規網址與官方 RSS。
//
// 輸入已帶頻道 ID 時仍會抓一次頻道頁：名稱只在頁面上，
// 而官方 RSS 會陣發性地整段 404，不能拿來當名稱的來源。
func (c *Client) GetChannel(ctx context.Context, input string) (*model.Channel, error) {
	target, err := ParseTarget(input)
	if err != nil {
		return nil, err
	}

	page := target.Handle
	if target.ID != "" {
		page = "channel/" + target.ID
	}
	html, err := c.fetch(ctx, fmt.Sprintf("%s/%s", c.baseURL, page))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("channel %s not found", input)
		}
		return nil, err
	}

	id := target.ID
	if id == "" {
		var ok bool
		if id, ok = ExtractChannelID(html); !ok {
			return nil, fmt.Errorf("no channel id found for %s", target.Handle)
		}
	}

	title, _ := ExtractChannelTitle(html)
	return &model.Channel{
		Platform: model.PLATFORM_YOUTUBE,
		ID:       id,
		Handle:   target.Handle,
		Title:    title,
		URL:      c.channelURL(id),
		RSS:      c.rssURL(id),
	}, nil
}

func (c *Client) channelURL(id string) string {
	return fmt.Sprintf("%s/channel/%s", c.baseURL, id)
}

func (c *Client) rssURL(id string) string {
	return fmt.Sprintf("%s/feeds/videos.xml?channel_id=%s", c.baseURL, id)
}

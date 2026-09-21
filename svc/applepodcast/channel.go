package applepodcast

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bizshuk/gx/model"
)

// lookupResponse 是 iTunes lookup 回應中本命令用到的欄位。
type lookupResponse struct {
	ResultCount int `json:"resultCount"`
	Results     []struct {
		Kind           string `json:"kind"`
		CollectionID   int64  `json:"collectionId"`
		CollectionName string `json:"collectionName"`
		FeedURL        string `json:"feedUrl"`
	} `json:"results"`
}

// GetChannel 由 collection ID 或節目網址解析出節目 ID、名稱、正規網址與 RSS。
//
// RSS 是 lookup 回傳的 feedUrl —— 節目發佈者自己的 feed，
// 地位等同 YouTube 的官方 RSS，不是第三方合成源。
func (c *Client) GetChannel(ctx context.Context, input string) (*model.Channel, error) {
	target, err := ParseTarget(input)
	if err != nil {
		return nil, err
	}

	q := url.Values{"id": {target.ID}, "entity": {"podcast"}}
	body, err := c.fetch(ctx, c.baseURL+"/lookup?"+q.Encode())
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("podcast %s not found", input)
		}
		return nil, err
	}

	var resp lookupResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse lookup for %s: %w", target.ID, err)
	}
	for _, r := range resp.Results {
		if r.Kind != "podcast" || strconv.FormatInt(r.CollectionID, 10) != target.ID {
			continue
		}
		return &model.Channel{
			Platform: model.PLATFORM_APPLE_PODCAST,
			ID:       target.ID,
			Title:    r.CollectionName,
			URL:      showURL(target.ID),
			RSS:      strings.TrimSpace(r.FeedURL),
		}, nil
	}
	return nil, fmt.Errorf("podcast %s not found", input)
}

// showURL 組出不帶地區與 slug 的正規節目頁；Apple 會依瀏覽者地區自行導向。
func showURL(id string) string {
	return fmt.Sprintf("%s/podcast/id%s", SHOW_BASE_URL, id)
}

// Package applepodcast 提供不需 API key 的 Apple Podcasts 節目查詢能力。
package applepodcast

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	gohttp "github.com/bizshuk/gosdk/http"
)

// ErrNotFound 代表 lookup 查無此節目（resultCount 為 0 或 HTTP 404）。
var ErrNotFound = errors.New("not found")

const (
	// DEFAULT_BASE_URL 是 iTunes Search API 的來源網域（lookup 端點所在）。
	DEFAULT_BASE_URL = "https://itunes.apple.com"
	// DEFAULT_TIMEOUT 單次 HTTP 請求的上限。
	DEFAULT_TIMEOUT = 10 * time.Second
	// DEFAULT_USER_AGENT 與其他領域一致，以桌面瀏覽器身分請求。
	DEFAULT_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36"
	// SHOW_BASE_URL 是節目頁的正規網域；lookup 與節目頁不同源。
	SHOW_BASE_URL = "https://podcasts.apple.com"
	// maxBodyBytes 限制讀入的回應大小；lookup 回應通常只有數 KB。
	maxBodyBytes = 1 << 20
)

// Client 是 iTunes lookup 的查詢器，重試策略沿用 gosdk 的預設預算。
type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	policy     gohttp.RetryPolicy
}

// NewClient 建立 Client，未指定的選項套用套件預設值。
func NewClient(opts ...Option) *Client {
	o := applyOptions(opts...)
	return &Client{
		httpClient: &http.Client{Timeout: o.timeout},
		baseURL:    strings.TrimRight(o.baseURL, "/"),
		userAgent:  o.userAgent,
		policy:     gohttp.DefaultRetryPolicy(),
	}
}

// fetch 取得指定網址的回應內容；429 與 5xx 視為暫時性失敗並重試。
func (c *Client) fetch(ctx context.Context, url string) ([]byte, error) {
	return gohttp.Retry(ctx, c.policy, func(ctx context.Context) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("build request %s: %w", url, err)
		}
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, gohttp.Retryable(fmt.Errorf("get %s: %w", url, err))
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("get %s: %w", url, ErrNotFound)
		}
		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("get %s: unexpected status %d", url, resp.StatusCode)
			if gohttp.IsRetryableStatus(resp.StatusCode) {
				return nil, gohttp.Retryable(err)
			}
			return nil, err
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
		if err != nil {
			return nil, gohttp.Retryable(fmt.Errorf("read %s: %w", url, err))
		}
		return body, nil
	})
}

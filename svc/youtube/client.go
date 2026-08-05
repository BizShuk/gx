// Package youtube 提供不需 API key 的 YouTube 公開頁面查詢能力。
package youtube

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	gohttp "github.com/bizshuk/gosdk/http"
)

// ErrNotFound 代表目標頁面不存在（HTTP 404）。
var ErrNotFound = errors.New("not found")

const (
	// DEFAULT_BASE_URL 是 YouTube 公開頁面的來源網域。
	DEFAULT_BASE_URL = "https://www.youtube.com"
	// DEFAULT_TIMEOUT 單次 HTTP 請求的上限。
	DEFAULT_TIMEOUT = 10 * time.Second
	// DEFAULT_USER_AGENT 以桌面瀏覽器身分請求，避免取得精簡版頁面。
	DEFAULT_USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) " +
		"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0 Safari/537.36"
	// maxBodyBytes 限制讀入的頁面大小；頻道 ID 出現在 HTML 前段，
	// 不需要為此把整份數 MB 的頁面收進記憶體。
	maxBodyBytes = 4 << 20
)

// Client 是 YouTube 公開頁面的抓取器，重試策略沿用 gosdk 的預設預算。
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
		baseURL:    o.baseURL,
		userAgent:  o.userAgent,
		policy:     gohttp.DefaultRetryPolicy(),
	}
}

// fetch 取得指定網址的 HTML 內容；429 與 5xx 視為暫時性失敗並重試。
func (c *Client) fetch(ctx context.Context, url string) (string, error) {
	return gohttp.Retry(ctx, c.policy, func(ctx context.Context) (string, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return "", fmt.Errorf("build request %s: %w", url, err)
		}
		req.Header.Set("User-Agent", c.userAgent)
		req.Header.Set("Accept-Language", "en-US,en;q=0.9")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return "", gohttp.Retryable(fmt.Errorf("get %s: %w", url, err))
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("get %s: %w", url, ErrNotFound)
		}
		if resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("get %s: unexpected status %d", url, resp.StatusCode)
			if gohttp.IsRetryableStatus(resp.StatusCode) {
				return "", gohttp.Retryable(err)
			}
			return "", err
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyBytes))
		if err != nil {
			return "", gohttp.Retryable(fmt.Errorf("read %s: %w", url, err))
		}
		return string(body), nil
	})
}

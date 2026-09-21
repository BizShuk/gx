// Package model 定義 gx 各平台共用的輸出物件。
package model

// 平台代碼，對應 Channel.Platform。
const (
	PLATFORM_YOUTUBE  = "youtube"
	PLATFORM_BILIBILI = "bilibili"
	// PLATFORM_APPLE_PODCAST 與命令的 domain 名一致。
	PLATFORM_APPLE_PODCAST = "apple-podcast"
)

// Channel 是所有平台「頻道」查詢的標準輸出物件。
//
// 欄位對所有平台一致，平台沒有的資訊留空並在輸出時省略 ——
// 例如 Bilibili 沒有官方 RSS，就不會有 rss，而不是代填第三方源。
// 新增欄位時只改這裡與 Fields，--json 與逐行輸出會一起長出來。
type Channel struct {
	Platform string `json:"platform"`
	ID       string `json:"id"`
	Handle   string `json:"handle,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"url"`
	RSS      string `json:"rss,omitempty"`
}

// Field 是逐行輸出的一個欄位。
type Field struct {
	Key   string
	Value string
}

// Fields 依固定順序列出非空欄位，key 與 JSON 欄位名一致。
func (c Channel) Fields() []Field {
	all := []Field{
		{"platform", c.Platform},
		{"id", c.ID},
		{"handle", c.Handle},
		{"title", c.Title},
		{"url", c.URL},
		{"rss", c.RSS},
	}
	fields := all[:0]
	for _, f := range all {
		if f.Value != "" {
			fields = append(fields, f)
		}
	}
	return fields
}

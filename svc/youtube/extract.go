package youtube

import "regexp"

// channelIDPattern 匹配頻道 ID 本身：UC 開頭加 22 個 base64url 字元。
var channelIDPattern = regexp.MustCompile(`^UC[A-Za-z0-9_-]{22}$`)

// channelIDInPage 依可靠度排序的頁面比對規則。
//
// 頻道頁的 HTML 裡散落著大量同樣長度、同樣字元集的隨機 token
// （visitor data、播放清單 id 等），只比對 ID 形狀會抓到錯的值，
// 因此一律要求 ID 出現在 channel 網址或 channelId 欄位的上下文中。
var channelIDInPage = []*regexp.Regexp{
	regexp.MustCompile(`<link\s+rel="canonical"\s+href="https?://[^"]*/channel/(UC[A-Za-z0-9_-]{22})"`),
	regexp.MustCompile(`"(?:channelId|externalId|externalChannelId)"\s*:\s*"(UC[A-Za-z0-9_-]{22})"`),
	regexp.MustCompile(`youtube\.com/channel/(UC[A-Za-z0-9_-]{22})`),
	regexp.MustCompile(`itemprop="identifier"\s+content="(UC[A-Za-z0-9_-]{22})"`),
}

// ExtractChannelID 從頁面 HTML 取出頻道 ID。
func ExtractChannelID(html string) (string, bool) {
	for _, pattern := range channelIDInPage {
		if match := pattern.FindStringSubmatch(html); match != nil {
			return match[1], true
		}
	}
	return "", false
}

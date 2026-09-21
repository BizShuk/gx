# gx — General eXtractor CLI

從公開網頁擷取結構化識別資訊的通用 CLI。以「不需要 API key、不需要帳號授權」為前提，
把平時散落在各處的一次性 `curl | grep` 收斂成可重複、可測試、可組合的子命令。

## 業務定義 (Business Definition)

日常工作中經常需要把一個「人看得懂的名字」換成「機器需要的 ID」——
YouTube 的 `@handle` 換成 `UCxxx` 頻道 ID 就是典型例子。這類查詢的共通樣貌是：

1. 對公開頁面發一次請求
2. 從回應中抽出唯一識別碼
3. 輸出成可直接餵給下一個工具的形式（純文字 / JSON）

`gx` 把這個樣貌標準化，讓每個新領域只需要補上「請求哪個頁面、抽哪個欄位」。

## 命令結構 (Command Structure)

命令採三段式 `gx <domain> <verb> <resource>`：

| 段落       | 意義           | 範例                     |
| ---------- | -------------- | ------------------------ |
| `domain`   | 資料來源領域   | `youtube`、`bilibili`、`apple-podcast` |
| `verb`     | 動作           | `get`                    |
| `resource` | 目標資源       | `channel`                |

## 領域流程 (Domain Flow)

### youtube get channel

把頻道的任意寫法解析成頻道 ID、名稱、正規網址與官方 RSS。
預設每行一個 `key: value`、多筆之間空一行；`--json` 一律輸出陣列。

```bash
gx youtube get channel @YouTube
# platform: youtube
# id: UCBR8-60-B28hp2BmDPdntcQ
# handle: @YouTube
# title: YouTube
# url: https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ
# rss: https://www.youtube.com/feeds/videos.xml?channel_id=UCBR8-60-B28hp2BmDPdntcQ

gx youtube get channel https://www.youtube.com/@YouTube/videos --json
# [
#   {
#     "platform": "youtube",
#     "id": "UCBR8-60-B28hp2BmDPdntcQ",
#     "handle": "@YouTube",
#     "title": "YouTube",
#     "url": "https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ",
#     "rss": "https://www.youtube.com/feeds/videos.xml?channel_id=UCBR8-60-B28hp2BmDPdntcQ"
#   }
# ]

cat handles.txt | gx youtube get channel --json | jq -r '.[].id'
```

接受的輸入寫法：

- `@YouTube`、`YouTube`
- `https://www.youtube.com/@YouTube`（含 `/videos` 等子頁與 `?si=` 追蹤參數）
- `youtube.com/c/YouTube`、`youtube.com/user/YouTube`
- `https://www.youtube.com/channel/UCxxx`、`UCxxx`（已含 ID 時仍抓頁面取名稱）

流程：正規化輸入 → 取得頻道頁 HTML → 在 canonical link / `channelId` 欄位的上下文中
比對 ID、取出頻道名稱 → 組出正規網址與官方 RSS（`/feeds/videos.xml?channel_id=`）。找不到頻道時以
`channel @xxx not found` 結束，離開碼 1。

### bilibili get channel

把 UP 主空間的 UID 或網址解析成正規的空間網址。Bilibili 沒有官方 channel RSS，
本命令不輸出 feed、也不代填 RSSHub 等第三方合成源。

```bash
gx bilibili get channel https://space.bilibili.com/2267573/video
# platform: bilibili
# id: 2267573
# url: https://space.bilibili.com/2267573

gx bilibili get channel 2267573 --json
# [
#   {
#     "platform": "bilibili",
#     "id": "2267573",
#     "url": "https://space.bilibili.com/2267573"
#   }
# ]
```

接受的輸入寫法：

- `2267573`（裸 UID）
- `https://space.bilibili.com/2267573`（含 `/video`、合集子頁與追蹤參數）
- `https://m.bilibili.com/space/2267573`

UID 已在輸入裡時不發請求。單支影片網址、暱稱、`b23.tv` 短鏈不是空間識別碼，以錯誤結束。

### apple-podcast get channel

把 Apple Podcasts 節目的 collection ID 或網址解析成節目 ID、名稱、正規網址與 RSS。
經 iTunes lookup 查詢，不需要 API key。`rss` 是節目發佈者自己的 feed（lookup 的 `feedUrl`），
地位等同 YouTube 的官方 RSS。

```bash
gx apple-podcast get channel 'https://podcasts.apple.com/tw/podcast/xxx/id1702409419?l=en-GB'
# platform: apple-podcast
# id: 1702409419
# title: 科技浪 Tech.wav
# url: https://podcasts.apple.com/podcast/id1702409419
# rss: https://feed.firstory.me/rss/user/cm3o5681s06e801v3fxpjehwb

gx apple-podcast get channel 1702409419 --json
```

接受的輸入寫法：

- `1702409419`、`id1702409419`
- `https://podcasts.apple.com/<地區>/podcast/<slug>/id1702409419`（含 `?l=`、單集 `?i=` 參數）
- `https://itunes.apple.com/us/podcast/id1702409419`

`url` 一律組成不帶地區與 slug 的 `https://podcasts.apple.com/podcast/id<id>`，由 Apple 依瀏覽者地區導向。
lookup 對不存在的 ID 回 200 加空結果，本命令將其視為 `podcast xxx not found`，離開碼 1。

## 設定 (Configuration)

設定檔位於 `~/.config/gx/settings.json`，首次執行自動建立。

| Key                | 預設值                     | 用途                       |
| ------------------ | -------------------------- | -------------------------- |
| `log_level`        | `info`                     | 日誌等級                   |
| `http_timeout`     | `10s`                      | 單次 HTTP 請求上限         |
| `http_user_agent`  | 桌面版 Chrome UA           | 避免拿到精簡版頁面         |
| `youtube_base_url` | `https://www.youtube.com`  | YouTube 頁面來源網域       |
| `bilibili_base_url` | `https://space.bilibili.com` | Bilibili 空間頁來源網域 |
| `apple_podcast_base_url` | `https://itunes.apple.com` | iTunes lookup 來源網域 |

檢視與修改：`gx config`（由 gosdk 提供）。

## 安裝 (Install)

```bash
npm run deploy   # go install .
```

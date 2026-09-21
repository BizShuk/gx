# 術語表 (Terminology)

## 命令結構

| 術語       | 定義                                                                 |
| ---------- | -------------------------------------------------------------------- |
| `domain`   | 資料來源領域，命令樹第一段，對應 `cmd/<domain>/` 與 `svc/<domain>/`。目前有 `youtube`、`bilibili`、`apple-podcast`。 |
| `verb`     | 動作，命令樹第二段。目前只有 `get`（唯讀查詢）。                     |
| `resource` | 目標資源，命令樹第三段，實際發出請求並輸出結果的葉命令。             |
| `Channel`  | 所有平台 `get channel` 共用的標準輸出物件：`platform`、`id`、`handle`、`title`、`url`、`rss`，平台沒有的欄位省略。 |
| 逐行輸出   | 未帶 `--json` 時的預設呈現：每行一個 `key: value`，多筆之間空一行。 |
| `--json`   | 一律輸出 `Channel` 陣列，單筆也是陣列。                              |

## YouTube 領域

| 術語                | 定義                                                                             |
| ------------------- | -------------------------------------------------------------------------------- |
| handle              | YouTube 頻道的人類可讀名稱，以 `@` 開頭，如 `@YouTube`。本專案一律正規化成含 `@`。 |
| channel ID          | 頻道的機器識別碼，`UC` 開頭加 22 個 base64url 字元，如 `UCBR8-60-B28hp2BmDPdntcQ`。 |
| canonical URL       | 由 channel ID 組成的正規頻道網址 `<base>/channel/<id>`，輸出的 `url` 欄位。         |
| official RSS        | YouTube 平台自己發的頻道 feed：`<base>/feeds/videos.xml?channel_id=<id>`。        |
| `Target`            | `svc/youtube` 對輸入的解析結果：帶 ID 時填 `ID`，否則填 `Handle`。               |
| lookalike token     | 頁面中與 channel ID 同形狀但意義不同的隨機字串（visitor data 等），必須排除。     |

## Bilibili 領域

| 術語          | 定義                                                                                         |
| ------------- | -------------------------------------------------------------------------------------------- |
| space UID     | UP 主空間的機器識別碼，十進位數字、不以 0 開頭，如 `2267573`。對應 URL 路徑上的那串數字。     |
| space URL     | 由 UID 組成的正規空間網址 `<base>/<uid>`，預設 `https://space.bilibili.com/<uid>`。           |
| `channel`     | 命令樹的 resource 名，與 youtube 對齊；在本領域對應的是 UP 主空間，不是合集或單支影片。       |
| official RSS  | 不存在。Bilibili 不為空間提供 feed；第三方合成源（RSSHub 等）不是本工具的輸出。               |

## Apple Podcasts 領域

| 術語            | 定義                                                                                   |
| --------------- | -------------------------------------------------------------------------------------- |
| collection ID   | 節目的機器識別碼，十進位數字，對應網址路徑上的 `id<數字>` 段，如 `1702409419`。          |
| iTunes lookup   | Apple 的公開 JSON 查詢端點 `<base>/lookup?id=<id>&entity=podcast`，不需要 API key。      |
| show URL        | 正規節目頁 `https://podcasts.apple.com/podcast/id<id>`，輸出的 `url` 欄位。             |
| feed URL        | lookup 回傳的 `feedUrl`，節目發佈者自己的 RSS，輸出的 `rss` 欄位。                      |
| `channel`       | 命令樹的 resource 名，與 youtube 對齊；在本領域對應的是節目 (show)，不是單集。          |

## 設定

| 術語         | 定義                                                                       |
| ------------ | -------------------------------------------------------------------------- |
| 扁平 key     | 不含巢狀結構的設定鍵，如 `http_timeout`，可被 `APP_HTTP_TIMEOUT` 環境變數覆寫。 |
| 設定根目錄   | `~/.config/gx/`，由 gosdk `config.Default(WithAppName("gx"))` 決定。        |

## 錯誤

| 術語          | 定義                                                        |
| ------------- | ----------------------------------------------------------- |
| `ErrNotFound` | 目標頁面回 404 的哨兵錯誤，由 service 層轉成領域語意的訊息。 |
| retryable     | 值得重試的暫時性失敗：連線錯誤、429、5xx。                  |

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
| `domain`   | 資料來源領域   | `youtube`                |
| `verb`     | 動作           | `get`                    |
| `resource` | 目標資源       | `channel`                |

## 領域流程 (Domain Flow)

### youtube get channel

把頻道的任意寫法解析成正規的 channel ID 與網址。

```bash
gx youtube get channel @YouTube
# https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ

gx youtube get channel @YouTube --id
# UCBR8-60-B28hp2BmDPdntcQ

gx youtube get channel https://www.youtube.com/@YouTube/videos --json
# {
#   "handle": "@YouTube",
#   "id": "UCBR8-60-B28hp2BmDPdntcQ",
#   "url": "https://www.youtube.com/channel/UCBR8-60-B28hp2BmDPdntcQ"
# }
```

接受的輸入寫法：

- `@YouTube`、`YouTube`
- `https://www.youtube.com/@YouTube`（含 `/videos` 等子頁與 `?si=` 追蹤參數）
- `youtube.com/c/YouTube`、`youtube.com/user/YouTube`
- `https://www.youtube.com/channel/UCxxx`、`UCxxx`（已含 ID 時不發請求）

流程：正規化輸入 → 取得頻道頁 HTML → 在 canonical link / `channelId` 欄位的上下文中
比對 ID → 組出正規網址。找不到頻道時以 `channel @xxx not found` 結束，離開碼 1。

## 設定 (Configuration)

設定檔位於 `~/.config/gx/settings.json`，首次執行自動建立。

| Key                | 預設值                     | 用途                       |
| ------------------ | -------------------------- | -------------------------- |
| `log_level`        | `info`                     | 日誌等級                   |
| `http_timeout`     | `10s`                      | 單次 HTTP 請求上限         |
| `http_user_agent`  | 桌面版 Chrome UA           | 避免拿到精簡版頁面         |
| `youtube_base_url` | `https://www.youtube.com`  | YouTube 頁面來源網域       |

檢視與修改：`gx config`（由 gosdk 提供）。

## 安裝 (Install)

```bash
npm run deploy   # go install .
```

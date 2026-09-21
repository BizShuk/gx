# gx — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

```tree
gx/
├── main.go                      # 進入點：載入設定後執行命令樹
├── go.mod                       # module github.com/bizshuk/gx, Go 1.26
├── cmd/                         # 命令層：只負責旗標解析與輸出，不含業務邏輯
│   ├── root.go                  # RootCmd + Execute()，掛載各領域命令
│   ├── youtube/                 # youtube 領域命令樹
│   │   ├── youtube.go           # Cmd：領域父命令
│   │   ├── get.go               # getCmd：讀取類動作
│   │   └── channel.go           # channelCmd：--json 旗標，流程交給 lookup
│   ├── bilibili/                # bilibili 領域命令樹（結構同 youtube）
│   ├── applepodcast/            # apple-podcast 領域命令樹（結構同 youtube）
│   └── lookup/
│       └── lookup.go            # 各平台 get channel 共用流程：讀目標、查詢、輸出
├── svc/                         # 服務層：對外請求與解析
│   ├── youtube/
│   │   ├── client.go            # Client、fetch()、ErrNotFound、預設常數
│   │   ├── option.go            # 函式選項 + viper key（設定讀取在此收斂）
│   │   ├── target.go            # ParseTarget()：輸入正規化成 handle 或 ID
│   │   ├── extract.go           # ExtractChannelID() / ExtractChannelTitle()：頁面比對規則
│   │   ├── channel.go           # GetChannel()：ID、名稱、官方 RSS
│   │   ├── target_test.go
│   │   └── channel_test.go
│   ├── bilibili/
│   │   ├── client.go            # Client、DEFAULT_BASE_URL（UID 已在輸入時不發請求）
│   │   ├── option.go            # 函式選項 + viper key
│   │   ├── target.go            # ParseTarget()：UID / 空間網址
│   │   └── channel.go           # GetChannel()；無 rss、無名稱
│   └── applepodcast/
│       ├── client.go            # Client、fetch()、DEFAULT_BASE_URL（iTunes lookup）
│       ├── option.go            # 函式選項 + viper key
│       ├── target.go            # ParseTarget()：collection ID / 節目網址
│       └── channel.go           # GetChannel()：lookup → 名稱、正規網址、feedUrl
├── model/
│   └── channel.go               # 各平台共用的標準輸出物件 Channel
├── render/
│   ├── json.go                  # --json：一律輸出陣列
│   └── lines.go                 # 預設：每行一個 key: value，記錄間空行
├── config/
│   ├── config.go                # config.Default()：gosdk 設定載入
│   └── default_settings.json    # 內嵌預設值，首次執行時寫入 ~/.config/gx/
├── docs/
│   ├── terminology.md           # 術語表
│   └── memory/                  # 決策與操作記錄
├── package.json                 # 統一任務入口（dev/test/build/deploy/lint）
├── README.md                    # 業務定義與命令用法
├── README.todo                  # 待辦
└── AGENTS.md -> CLAUDE.md
```

## 技術棧 (Tech Stack)

- Go 1.26
- `github.com/bizshuk/gosdk` v1.3.1 — 設定載入 (`config`)、HTTP 重試 (`http`)、`gx config` 子命令 (`cmd`)
- `spf13/cobra` — CLI 命令樹
- `spf13/viper` — 設定讀取

## 關鍵決策 (Key Decisions)

### 命令採 `<domain> <verb> <resource>` 三段式

新增領域 = 在 `cmd/` 下新增一個領域套件（父命令 + 動作 + 資源三個檔案），
再於 `cmd/root.go` 掛上去。`svc/` 下對應新增同名領域套件。
`render/` 是跨領域共用的輸出層，不放任何領域知識。

### 設定值由 svc 層自行讀 viper，不經由全域 struct 轉手

`svc/youtube/option.go` 的 `applyOptions()` 取值順序為
**函式選項 > viper 設定 > 套件預設常數**。這與 gosdk `db.InitSQLite()` 自行讀
`SQLITE_PATH` 的作法一致。`config` 套件只負責把設定載進 viper，
不再維護一份把設定欄位一路搬進 service 的樣板程式碼。

函式選項的存在理由是測試：`WithBaseURL(srv.URL)` 讓測試指向 `httptest` 伺服器。

### 輸出物件標準化：一個 Channel，兩種呈現

所有平台的 `get channel` 都回 `model.Channel`（`platform` / `id` / `handle` /
`title` / `url` / `rss`），平台沒有的欄位留空並省略。輸出只有兩種：
預設逐行 `key: value`（多筆以空行分隔），`--json` 一律輸出陣列（單筆也是）。

不再提供 `--id`、`--rss` 這類單欄位旗標：每多一個欄位就要多一個旗標，
而且單筆物件、多筆陣列的切換讓呼叫端得先數目標才知道怎麼解析。
程式呼叫端一律 `--json` 後自取欄位。

### 頻道名稱來自頻道頁，不來自 RSS

名稱依序取 `og:title`（HTML unescape）與 `channelMetadataRenderer.title`（JSON 字串）。
官方 RSS 會對個別頻道持續回 404/500，不能當名稱來源。因此輸入已是頻道 ID 時
仍會抓一次 `/channel/<id>` 頁面。取不到名稱不算失敗，`title` 留空。

### Bilibili 沒有官方 channel RSS，不代填第三方源

YouTube 頻道有平台自己發的 `/feeds/videos.xml?channel_id=`，所以
`gx youtube get channel` 的輸出帶 `rss`。

Bilibili 的 UP 主空間沒有對等的官方 feed（舊的分區 `/rss-N.xml` 已下線）。
生態裡看得到的「B 站 RSS」是 RSSHub 等第三方打簽名 JSON API 再包成 XML。
那些網址不是來源站的公開頁、公共實例也不穩定，因此
`gx bilibili get channel` 只回 UID 與空間網址，輸出**沒有** `rss` 欄位。

### Apple Podcasts 走 iTunes lookup，不爬節目頁

節目頁是 JS 渲染的，名稱與 feed 都不在穩定的 HTML 位置；iTunes lookup
(`/lookup?id=<id>&entity=podcast`) 是公開 JSON，一次給齊 `collectionName` 與 `feedUrl`。
lookup 對不存在的 ID 回 `200` 加 `resultCount: 0`，不是 404 —— 空結果要自己轉成 not found。
domain 名 `apple-podcast` 帶平台前綴：podcast 是媒體型態而非平台，
日後的其他 podcast 目錄（Spotify 等）各自是一個 domain。

### 頻道 ID 一律在上下文中比對，不比對 ID 形狀

YouTube 頻道頁的 HTML 裡散落大量與頻道 ID 同形狀（`UC` + 22 字元 base64url）的隨機
token——visitor data、播放清單 id 等。只用 `UC[A-Za-z0-9_-]{22}` 掃全頁會抓到錯的值，
而且每次請求結果不同。因此 `extract.go` 依可靠度排序，要求 ID 出現在
canonical link、`channelId`/`externalId` 欄位、channel 網址或 `itemprop="identifier"`
的上下文中。**修改該檔時務必保留這個約束。**

### 錯誤輸出統一在 Execute

`RootCmd` 設 `SilenceErrors` 與 `SilenceUsage`，錯誤由 `Execute()` 印到 stderr 後
以離開碼 1 結束，避免 cobra 重複輸出同一則訊息。

### 重試策略沿用 gosdk 預算

`svc/youtube/client.go` 用 `gohttp.Retry` + `DefaultRetryPolicy()`（5 次、指數退避）。
429 與 5xx 標記為 `Retryable`；404 轉成 `ErrNotFound` 由 `GetChannel` 包成
`channel @xxx not found`；其餘 4xx 直接失敗，不消耗重試預算。

## 開發 (Development)

```bash
npm run dev -- youtube get channel @YouTube   # go run .
npm test                                       # go test ./...
npm run lint                                   # gofmt -l . && go vet ./...
npm run build                                  # go build -o bin/gx .
```

測試不打外部網路：`svc/youtube` 的測試全部走 `httptest` 伺服器，
以 `WithBaseURL` 注入。`svc/bilibili` 的 UID 已在輸入裡，測試不發 HTTP。
`svc/applepodcast` 與 `cmd/applepodcast` 同樣以 `httptest` 假冒 lookup。

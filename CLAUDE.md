# gx — 技術脈絡 (Technical Context)

## 專案結構 (Project Structure)

```tree
gx/
├── main.go                      # 進入點：載入設定後執行命令樹
├── go.mod                       # module github.com/bizshuk/gx, Go 1.26
├── cmd/                         # 命令層：只負責旗標解析與輸出，不含業務邏輯
│   ├── root.go                  # RootCmd + Execute()，掛載各領域命令
│   └── youtube/                 # youtube 領域命令樹
│       ├── youtube.go           # Cmd：領域父命令
│       ├── get.go               # getCmd：讀取類動作
│       └── channel.go           # channelCmd：--json / --id 旗標與輸出
├── svc/                         # 服務層：對外請求與解析
│   └── youtube/
│       ├── client.go            # Client、fetch()、ErrNotFound、預設常數
│       ├── option.go            # 函式選項 + viper key（設定讀取在此收斂）
│       ├── target.go            # ParseTarget()：輸入正規化成 handle 或 ID
│       ├── extract.go           # ExtractChannelID()：頁面 HTML 的 ID 比對規則
│       ├── channel.go           # Channel 模型 + GetChannel()
│       ├── target_test.go
│       └── channel_test.go
├── render/
│   └── json.go                  # 各子命令共用的 --json 輸出
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
以 `WithBaseURL` 注入。

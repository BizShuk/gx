// Package get 是不分平台的讀取命令：由網址網域自動判斷平台。
//
// 它是三段式 `<domain> <verb> <resource>` 省略 domain 的形式 ——
// domain 由輸入的網址決定，其餘旗標與輸出和各平台的同名命令一致。
package get

import (
	"github.com/bizshuk/gx/cmd/lookup"
	"github.com/bizshuk/gx/svc/resolve"
	"github.com/spf13/cobra"
)

// Cmd 匯集不分平台的讀取類命令。
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "依網址自動判斷平台並讀取資源",
}

var channelAsJSON bool

// channelCmd 依網址網域分派給對應平台的 get channel。
var channelCmd = &cobra.Command{
	Use:   "channel [url...]",
	Short: "依網址網域自動判斷平台並解析頻道",
	Long: `依網址網域自動判斷平台（youtube、bilibili、apple-podcast）並解析頻道，
輸出與 gx <domain> get channel 相同。只接受網址：@handle、裸 ID 說不出平台，
請改用對應平台的命令。
預設每行一個 key: value 欄位、多筆之間以空行分隔；--json 一律輸出 JSON 陣列。
支援傳入多個目標或由 stdin 讀取，同一批可混合不同平台。

範例:
  gx get channel https://www.youtube.com/@YouTube
  gx get channel 'https://podcasts.apple.com/tw/podcast/xxx/id1702409419' --json
  cat urls.txt | gx get channel --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lookup.Run(cmd, args, resolve.New().Channel, channelAsJSON)
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 陣列輸出完整結果")
	Cmd.AddCommand(channelCmd)
}

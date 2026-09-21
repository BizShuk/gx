package applepodcast

import (
	"github.com/bizshuk/gx/cmd/lookup"
	svc "github.com/bizshuk/gx/svc/applepodcast"
	"github.com/spf13/cobra"
)

var channelAsJSON bool

// channelCmd 由節目網址解析出節目資訊。
var channelCmd = &cobra.Command{
	Use:   "channel [url...]",
	Short: "解析節目 ID、名稱、正規網址與 RSS",
	Long: `由 Apple Podcasts 節目網址解析出節目 ID、名稱、正規網址與發佈者的 RSS feed。
只接受 podcasts.apple.com／itunes.apple.com 網址，不收裸 collection ID。
預設每行一個 key: value 欄位、多筆之間以空行分隔；--json 一律輸出 JSON 陣列。
支援傳入多個目標或由 stdin 讀取。

範例:
  gx apple-podcast get channel 'https://podcasts.apple.com/tw/podcast/xxx/id1702409419'
  gx apple-podcast get channel 'https://podcasts.apple.com/podcast/id1702409419' --json
  cat urls.txt | gx apple-podcast get channel --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lookup.Run(cmd, args, svc.NewClient().GetChannel, channelAsJSON)
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 陣列輸出完整結果")
}

package youtube

import (
	"github.com/bizshuk/gx/cmd/lookup"
	svc "github.com/bizshuk/gx/svc/youtube"
	"github.com/spf13/cobra"
)

var channelAsJSON bool

// channelCmd 由 handle 解析頻道 ID、名稱、正規網址與官方 RSS。
var channelCmd = &cobra.Command{
	Use:   "channel [handle|url|id...]",
	Short: "解析頻道 ID、名稱、正規網址與官方 RSS",
	Long: `由 @handle、頻道網址或頻道 ID 解析出頻道 ID、名稱、正規的 channel 網址與官方 RSS feed。
預設每行一個 key: value 欄位、多筆之間以空行分隔；--json 一律輸出 JSON 陣列。
支援傳入多個目標或由 stdin 讀取。

範例:
  gx youtube get channel @YouTube
  gx youtube get channel @YouTube @NASA --json
  gx youtube get channel @YouTube --json | jq -r '.[0].title'
  gx youtube get channel - < handles.txt`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lookup.Run(cmd, args, svc.NewClient().GetChannel, channelAsJSON)
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 陣列輸出完整結果")
}

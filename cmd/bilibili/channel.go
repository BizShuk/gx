package bilibili

import (
	"github.com/bizshuk/gx/cmd/lookup"
	svc "github.com/bizshuk/gx/svc/bilibili"
	"github.com/spf13/cobra"
)

var channelAsJSON bool

// channelCmd 由 UID 或空間網址解析出正規的空間網址。
var channelCmd = &cobra.Command{
	Use:   "channel [uid|url...]",
	Short: "解析 UP 主空間 UID 與正規網址",
	Long: `由空間 UID 或 space.bilibili.com 網址解析出正規的空間網址。
Bilibili 沒有官方 channel RSS，本命令不輸出 feed。
預設每行一個 key: value 欄位、多筆之間以空行分隔；--json 一律輸出 JSON 陣列。
支援傳入多個目標或由 stdin 讀取。

範例:
  gx bilibili get channel 2267573
  gx bilibili get channel https://space.bilibili.com/2267573/video --json
  cat uids.txt | gx bilibili get channel --json`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return lookup.Run(cmd, args, svc.NewClient().GetChannel, channelAsJSON)
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 陣列輸出完整結果")
}

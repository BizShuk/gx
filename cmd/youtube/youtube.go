// Package youtube 定義 gx 的 youtube 領域命令樹。
package youtube

import "github.com/spf13/cobra"

// Cmd 是 youtube 領域的父命令。
var Cmd = &cobra.Command{
	Use:   "youtube",
	Short: "YouTube 公開資訊查詢",
	Long:  "從 YouTube 公開頁面擷取資訊，不需要 API key。",
}

func init() {
	Cmd.AddCommand(getCmd)
}

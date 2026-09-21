// Package applepodcast 定義 gx 的 apple-podcast 領域命令樹。
package applepodcast

import "github.com/spf13/cobra"

// Cmd 是 apple-podcast 領域的父命令。
var Cmd = &cobra.Command{
	Use:   "apple-podcast",
	Short: "Apple Podcasts 公開資訊查詢",
	Long:  "經 iTunes lookup 查詢 Apple Podcasts 節目資訊，不需要 API key。",
}

func init() {
	Cmd.AddCommand(getCmd)
}

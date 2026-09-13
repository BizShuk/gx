// Package bilibili 定義 gx 的 bilibili 領域命令樹。
package bilibili

import "github.com/spf13/cobra"

// Cmd 是 bilibili 領域的父命令。
var Cmd = &cobra.Command{
	Use:   "bilibili",
	Short: "Bilibili 公開資訊查詢",
	Long:  "從 Bilibili 公開網址擷取資訊，不需要 API key。平台沒有官方 channel RSS。",
}

func init() {
	Cmd.AddCommand(getCmd)
}

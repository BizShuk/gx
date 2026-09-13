package bilibili

import "github.com/spf13/cobra"

// getCmd 匯集 bilibili 領域的讀取類命令。
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "讀取 Bilibili 資源",
}

func init() {
	getCmd.AddCommand(channelCmd)
}

package youtube

import "github.com/spf13/cobra"

// getCmd 匯集 youtube 領域的讀取類命令。
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "讀取 YouTube 資源",
}

func init() {
	getCmd.AddCommand(channelCmd)
}

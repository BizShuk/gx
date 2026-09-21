package applepodcast

import "github.com/spf13/cobra"

// getCmd 匯集 apple-podcast 領域的讀取類命令。
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "讀取 Apple Podcasts 資源",
}

func init() {
	getCmd.AddCommand(channelCmd)
}

// Package cmd 組裝 gx 的命令樹：根命令只負責掛載各領域子命令。
package cmd

import (
	"fmt"
	"os"

	gosdkcmd "github.com/bizshuk/gosdk/cmd"
	"github.com/bizshuk/gx/cmd/applepodcast"
	"github.com/bizshuk/gx/cmd/bilibili"
	"github.com/bizshuk/gx/cmd/youtube"
	"github.com/spf13/cobra"
)

// RootCmd 是 CLI 的根命令。
var RootCmd = &cobra.Command{
	Use:   "gx",
	Short: "General eXtractor CLI",
	Long: `gx 是一組通用擷取子命令的集合。

命令採 <domain> <verb> <resource> 三段式，新增領域時只要在 cmd/ 下
新增一個領域套件並掛到 RootCmd 即可。`,
	// 錯誤一律由 Execute 統一輸出到 stderr，避免 cobra 再印一次。
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute 執行根命令。
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(gosdkcmd.ConfigCmd, youtube.Cmd, bilibili.Cmd, applepodcast.Cmd)
}

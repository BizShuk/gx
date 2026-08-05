package youtube

import (
	"fmt"

	"github.com/bizshuk/gx/render"
	svc "github.com/bizshuk/gx/svc/youtube"
	"github.com/spf13/cobra"
)

var (
	channelAsJSON bool
	channelIDOnly bool
)

// channelCmd 由 handle 解析頻道 ID 與正規網址。
var channelCmd = &cobra.Command{
	Use:   "channel <handle|url|id>",
	Short: "解析頻道 ID 與正規網址",
	Long: `由 @handle、頻道網址或頻道 ID 解析出正規的 channel 網址。

範例:
  gx youtube get channel @YouTube
  gx youtube get channel https://www.youtube.com/@YouTube --json
  gx youtube get channel YouTube --id`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		channel, err := svc.NewClient().GetChannel(cmd.Context(), args[0])
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		switch {
		case channelAsJSON:
			return render.JSON(out, channel)
		case channelIDOnly:
			_, err = fmt.Fprintln(out, channel.ID)
		default:
			_, err = fmt.Fprintln(out, channel.URL)
		}
		return err
	},
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 輸出完整結果")
	channelCmd.Flags().BoolVar(&channelIDOnly, "id", false, "只輸出頻道 ID")
}

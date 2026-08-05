package youtube

import (
	"bufio"
	"fmt"
	"strings"

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
	Use:   "channel [handle|url|id...]",
	Short: "解析頻道 ID 與正規網址",
	Long: `由 @handle、頻道網址或頻道 ID 解析出正規的 channel 網址。
支援傳入多個目標或由 stdin 讀取。

範例:
  gx youtube get channel @YouTube
  gx youtube get channel @YouTube @NASA --json
  cat handles.txt | gx youtube get channel --id
  gx youtube get channel - < handles.txt`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		targets, err := parseTargets(cmd, args)
		if err != nil {
			return err
		}

		client := svc.NewClient()
		out := cmd.OutOrStdout()
		errOut := cmd.ErrOrStderr()

		var results []*svc.Channel
		var hasError bool

		for _, target := range targets {
			ch, err := client.GetChannel(cmd.Context(), target)
			if err != nil {
				hasError = true
				fmt.Fprintf(errOut, "error fetching %s: %v\n", target, err)
				continue
			}
			results = append(results, ch)
		}

		if len(results) == 0 && hasError {
			return fmt.Errorf("failed to resolve channel(s)")
		}

		if channelAsJSON {
			if len(targets) == 1 {
				if len(results) == 1 {
					if err := render.JSON(out, results[0]); err != nil {
						return err
					}
				}
			} else {
				if err := render.JSON(out, results); err != nil {
					return err
				}
			}
		} else {
			for _, ch := range results {
				var line string
				if channelIDOnly {
					line = ch.ID
				} else {
					line = ch.URL
				}
				if _, err := fmt.Fprintln(out, line); err != nil {
					return err
				}
			}
		}

		if hasError {
			return fmt.Errorf("one or more channels failed to resolve")
		}
		return nil
	},
}

func parseTargets(cmd *cobra.Command, args []string) ([]string, error) {
	if len(args) > 0 && !(len(args) == 1 && args[0] == "-") {
		return args, nil
	}

	var targets []string
	scanner := bufio.NewScanner(cmd.InOrStdin())
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			targets = append(targets, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read stdin: %w", err)
	}
	if len(targets) == 0 {
		return nil, fmt.Errorf("at least one target (handle, url, or id) is required")
	}
	return targets, nil
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 輸出完整結果")
	channelCmd.Flags().BoolVar(&channelIDOnly, "id", false, "只輸出頻道 ID")
}

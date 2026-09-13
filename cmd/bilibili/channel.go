package bilibili

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bizshuk/gx/render"
	svc "github.com/bizshuk/gx/svc/bilibili"
	"github.com/spf13/cobra"
)

var (
	channelAsJSON bool
	channelIDOnly bool
)

// channelCmd 由 UID 或空間網址解析出正規的空間網址。
var channelCmd = &cobra.Command{
	Use:   "channel [uid|url...]",
	Short: "解析 UP 主空間 UID 與正規網址",
	Long: `由空間 UID 或 space.bilibili.com 網址解析出正規的空間網址。
Bilibili 沒有官方 channel RSS，本命令不輸出 feed。
支援傳入多個目標或由 stdin 讀取。

範例:
  gx bilibili get channel 2267573
  gx bilibili get channel https://space.bilibili.com/2267573/video --id
  gx bilibili get channel 2267573 --json
  cat uids.txt | gx bilibili get channel --id`,
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
				line := ch.URL
				if channelIDOnly {
					line = ch.ID
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
		return nil, fmt.Errorf("at least one target (uid or url) is required")
	}
	return targets, nil
}

func init() {
	channelCmd.Flags().BoolVar(&channelAsJSON, "json", false, "以 JSON 輸出完整結果")
	channelCmd.Flags().BoolVar(&channelIDOnly, "id", false, "只輸出空間 UID")
}

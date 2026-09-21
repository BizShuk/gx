// Package lookup 是各平台 `get channel` 共用的命令流程：
// 讀目標、逐一查詢、依旗標輸出標準 Channel 物件。
package lookup

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"github.com/bizshuk/gx/model"
	"github.com/bizshuk/gx/render"
	"github.com/spf13/cobra"
)

// Resolver 把一個目標解析成標準 Channel。
type Resolver func(ctx context.Context, target string) (*model.Channel, error)

// Run 查詢所有目標並輸出。asJSON 時一律輸出 JSON 陣列，否則逐行輸出。
// 部分失敗時仍輸出成功的那幾筆，錯誤寫到 stderr 並以非零結束。
func Run(cmd *cobra.Command, args []string, resolve Resolver, asJSON bool) error {
	targets, err := targets(cmd, args)
	if err != nil {
		return err
	}

	errOut := cmd.ErrOrStderr()
	results := make([]*model.Channel, 0, len(targets))
	failed := 0
	for _, target := range targets {
		ch, err := resolve(cmd.Context(), target)
		if err != nil {
			failed++
			fmt.Fprintf(errOut, "error fetching %s: %v\n", target, err)
			continue
		}
		results = append(results, ch)
	}

	if len(results) == 0 {
		return fmt.Errorf("failed to resolve channel(s)")
	}

	out := cmd.OutOrStdout()
	if asJSON {
		err = render.JSON(out, results)
	} else {
		err = render.Lines(out, results)
	}
	if err != nil {
		return err
	}

	if failed > 0 {
		return fmt.Errorf("%d of %d channel(s) failed to resolve", failed, len(targets))
	}
	return nil
}

// targets 取命令列參數；沒有參數或只有 `-` 時改由 stdin 逐行讀取。
func targets(cmd *cobra.Command, args []string) ([]string, error) {
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
		return nil, fmt.Errorf("at least one target is required")
	}
	return targets, nil
}

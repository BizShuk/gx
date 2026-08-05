// Package render 統一各子命令的輸出格式。
package render

import (
	"encoding/json"
	"io"
)

// JSON 以縮排 JSON 輸出任意結果，供各子命令的 --json 旗標共用。
func JSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

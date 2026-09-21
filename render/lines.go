package render

import (
	"fmt"
	"io"

	"github.com/bizshuk/gx/model"
)

// Lines 以每行一個 `key: value` 輸出結果，記錄之間以空行分隔。
//
// 一個欄位一行而不是一筆一行的欄位清單：新增欄位不會挪動既有欄位的位置，
// 需要穩定結構的呼叫端改用 --json。
func Lines(w io.Writer, channels []*model.Channel) error {
	for i, ch := range channels {
		if i > 0 {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		for _, f := range ch.Fields() {
			if _, err := fmt.Fprintf(w, "%s: %s\n", f.Key, f.Value); err != nil {
				return err
			}
		}
	}
	return nil
}

package cli

import (
	"fmt"
	"io"
)

// PrintError 打印错误信息到指定的 io.Writer
func PrintError(w io.Writer, err error) {
	fmt.Fprintln(w, err)
}

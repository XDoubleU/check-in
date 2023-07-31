package tests

import (
	"fmt"
	"io"
)

func PrintJSON(body io.Reader) {
	b, _ := io.ReadAll(body)
	fmt.Println(string(b)) //nolint:forbidigo //intended println
}

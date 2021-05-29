package gradualtype

import (
	"fmt"
	"time"

	"github.com/oxide-one/systemd.go/pkg/color"
)

func GradualType(printStr string, delay int, colors ...string) {
	for _, color := range colors {
		fmt.Print(color)
	}
	timeDelay := time.Millisecond * time.Duration(delay)
	for _, letter := range printStr {
		fmt.Print(string(letter))
		time.Sleep(timeDelay)
	}
	fmt.Print(color.Reset)
	fmt.Println()
}

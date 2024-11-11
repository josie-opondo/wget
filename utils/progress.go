package utils

import (
	"fmt"
	"strings"
)

func ShowProgress(progress int64, total int64) {
	const length = 50
	if total <= 0 {
		fmt.Printf("\rDownloading...")
		return
	}
	percent := float64(progress) / float64(total) * 100
	numBars := int((percent / 100) * length)
	out := fmt.Sprintf("%.2f KiB / %.2f KiB [%s%s] %.0f%%",float64(progress)/1024, float64(total)/1024, strings.Repeat("=", numBars), strings.Repeat(" ", length-numBars), percent)
	fmt.Printf("\r%s", out)
}
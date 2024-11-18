package utils

import (
	"fmt"
	"math"
	"strings"
)

func roundToNearest(value float64) float64 {
	return math.Round(value*100) / 100
}

func size(val float64) string {
	if val > 1000000000 {
		res_gb := float64(val) / 1000000000
		return fmt.Sprintf("~%.2fGB", roundToNearest(res_gb))
	} else if val > 1000000 {
		res_gb := float64(val) / 1000000
		return fmt.Sprintf("[~%.2fMB]", roundToNearest(res_gb))
	}
	return fmt.Sprintf("%.0fKiB", val)
}

func ShowProgress(progress int64, total int64) {
	const length = 50
	if total <= 0 {
		fmt.Printf("\rDownloading...")
		return
	}
	percent := float64(progress) / float64(total) * 100
	numBars := int((percent / 100) * length)
	out := fmt.Sprintf("%s / %s [%s%s] %.0f%%", size(float64(progress)), size(float64(total)), strings.Repeat("=", numBars), strings.Repeat(" ", length-numBars), percent)
	fmt.Printf("\r%s", out)
}

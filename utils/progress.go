package utils

import (
	"fmt"
	"math"
	"strings"
	"time"
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

func (w WgetValues) ShowProgress(progress int64, total int64, startTime time.Time) {
	const length = 50
	if total <= 0 && !w.BackgroudMode {
		fmt.Printf("\rDownloading...")
		return
	}
	percent := float64(progress) / float64(total) * 100
	numBars := int((percent / 100) * length)

	// Calculate speed (bytes per second)
	elapsed := time.Since(startTime).Seconds()
	speed := float64(progress) / elapsed

	// Calculate estimated time remaining
	var eta string
	if speed > 0 {
		remaining := float64(total-progress) / speed
		eta = fmt.Sprintf("%02d:%02d:%02d", int(remaining/3600), int(remaining/60)%60, int(remaining)%60)
	} else {
		eta = "--:--:--"
	}

	// print the output if background mode is false
	if !w.BackgroudMode {
		out := fmt.Sprintf("%.2f KiB / %.2f KiB [%s%s] %.0f%% %s/s %s",
			float64(progress)/1024, float64(total)/1024,
			strings.Repeat("=", numBars), strings.Repeat(" ", length-numBars),
			percent, size(speed/1024), eta)
		fmt.Printf("\r%s", out)
	}
}

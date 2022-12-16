package tools

import "fmt"

func FormatNetSpeed(b float64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%.1f B/s", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB/s",
		float64(b)/float64(div), "KMGTPE"[exp])
}

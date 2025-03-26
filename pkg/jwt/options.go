package jwt

import (
	"strings"
	"time"
)

func ParseDuration(s string) time.Duration {
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		if d, err := time.ParseDuration(days + "h"); err == nil {
			return d * 24
		}
	}
	d, _ := time.ParseDuration(s)
	return d
}

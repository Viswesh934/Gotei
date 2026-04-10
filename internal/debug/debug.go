package debug

import (
	"log"
	"os"
	"strings"
)

var enabled = computeEnabled()

func computeEnabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("GOTEI_DEBUG")))
	switch v {
	case "1", "true", "yes", "on", "debug":
		return true
	default:
		return false
	}
}

func Enabled() bool {
	return enabled
}

func Logf(format string, args ...any) {
	if !enabled {
		return
	}
	log.Printf("[gotei-debug] "+format, args...)
}

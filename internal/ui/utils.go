package ui

import (
	"fmt"
	"strings"
	"time"
)

func RightPadTrim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func Trim(s string, length int) string {
	if len(s) >= length {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}
	return s
}

func HumanizeDuration(durationInSecs int) string {
	duration := time.Duration(durationInSecs) * time.Second

	if duration.Hours() > 48 {
		return fmt.Sprintf("%dd", int(duration.Hours()/24))
	}

	if duration.Seconds() < 60 {
		return fmt.Sprintf("%ds", int(duration.Seconds()))
	}

	if duration.Minutes() < 60 {
		return fmt.Sprintf("%dm", int(duration.Minutes()))
	}

	return fmt.Sprintf("%dh", int(duration.Hours()))
}

func allEqual(versions []versionInfo) bool {
	if len(versions) <= 1 {
		return true
	}

	for _, v := range versions {
		if v.errMsg != "" || v.notFound {
			return false
		}
	}

	versionsMap := make(map[string]struct{})

	for _, v := range versions {
		if v.errMsg != "" || v.notFound {
			return false
		}

		if v.version != "" {
			versionsMap[v.version] = struct{}{}
		}
	}

	return len(versionsMap) == 1
}

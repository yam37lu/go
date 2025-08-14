package utils

import "strings"

func SqlLikeEscape(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(key, "_", "\\_"), "%", "\\%")
}

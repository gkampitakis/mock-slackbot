package utils

import "regexp"

var tagsRegex = regexp.MustCompile(`(?i)<[!|@][\d\w]+>`)

func EscapeSlackTags(message string) string {
	return tagsRegex.ReplaceAllLiteralString(message, "")
}

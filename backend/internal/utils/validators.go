package utils

import (
	"regexp"
	"strings"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	tickerRegex   = regexp.MustCompile(`^[a-zA-Z]+$`)
)

func ValidateUsername(username string) bool {
	return len(username) >= 3 && len(username) <= 20 && usernameRegex.MatchString(username)
}

func ValidateTicker(ticker string) bool {
	return len(ticker) >= 2 && len(ticker) <= 15 && tickerRegex.MatchString(ticker)
}

func SanitizeTicker(firstName string) string {
	ticker := strings.TrimSpace(firstName)
	ticker = strings.ToUpper(ticker)
	return ticker
}

package opener

import (
	"errors"
	"strings"
)

// Target normalizes a resource value before passing it to a platform opener.
// Markdown links use the destination inside the final pair of parentheses;
// other non-empty values remain unchanged for compatibility.
func Target(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("resource value is empty")
	}
	if !strings.HasPrefix(value, "[") {
		return value, nil
	}

	separator := strings.Index(value, "](")
	if separator < 0 || !strings.HasSuffix(value, ")") {
		return "", errors.New("resource value is not a valid Markdown link")
	}
	target := strings.TrimSpace(value[separator+2 : len(value)-1])
	if target == "" || strings.ContainsAny(target, "\r\n") {
		return "", errors.New("resource value contains an empty Markdown link target")
	}
	return target, nil
}

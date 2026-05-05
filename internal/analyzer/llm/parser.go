package llm

import (
	"regexp"
	"strings"
)

func ExtractCommand(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "kubectl") {
			return strings.TrimSpace(line)
		}
	}

	re := regexp.MustCompile("`([^`]+)`")
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}

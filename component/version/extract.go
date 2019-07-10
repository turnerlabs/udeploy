package version

import (
	"fmt"
	"regexp"
)

// FormatExtract ...
func FormatExtract(image, regex string) string {

	version, build := Extract(image, regex)

	if len(build) > 0 {
		return fmt.Sprintf("%s.%s", version, build)
	}

	return version
}

// Extract ...
func Extract(image, regex string) (string, string) {
	tag := regexp.MustCompile(regex)

	matches := tag.FindAllStringSubmatch(image, -1)

	if missingVersion(matches) {
		return "", ""
	}

	if missingBuildNumber(matches) {
		return matches[0][1], ""
	}

	return matches[0][1], matches[0][2]
}

func missingVersion(matches [][]string) bool {
	if matches == nil || len(matches) < 1 || len(matches[0]) < 2 {
		return true
	}

	return false
}

func missingBuildNumber(matches [][]string) bool {
	if len(matches[0]) < 3 {
		return true
	}

	return false
}

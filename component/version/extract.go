package version

import (
	"fmt"
	"regexp"
)

// FormatExtract ...
func FormatExtract(image, regex string) (string, error) {

	version, build, err := Extract(image, regex)

	if len(build) > 0 {
		return fmt.Sprintf("%s.%s", version, build), err
	}

	return version, err
}

// Extract ...
func Extract(field, regex string) (string, string, error) {
	tag, err := regexp.Compile(regex)
	if err != nil {
		return "", "", err
	}

	matches := tag.FindAllStringSubmatch(field, -1)

	if missingVersion(matches) {
		return "", "", nil
	}

	if missingBuildNumber(matches) {
		return matches[0][1], "", nil
	}

	return matches[0][1], matches[0][2], nil
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

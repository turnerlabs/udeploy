package commit

import (
	"regexp"
	"strings"
)

const existingValueToken = "{EXISTING_VALUE}"

func mutate(text, existingValue, newValue string) string {
	var re = regexp.MustCompile(existingValue)

	for _, match := range re.FindAllString(text, -1) {
		newValue = strings.Replace(newValue, existingValueToken, match, -1)

		text = strings.Replace(text, match, newValue, -1)
	}

	text = strings.Replace(text, "\n", "<br />", -1)

	return text
}

package version

import (
	"fmt"
	"regexp"
)

// Undetermined ...
const Undetermined = "undetermined"

// Version ...
type Version struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	Divider string `json:"divider"`
}

// Copy ...
func (v Version) Copy() Version {
	return Version{
		Version: v.Version,
		Build:   v.Build,
		Divider: v.Divider,
	}
}

// Full ...
func (v Version) Full() string {
	if len(v.Version) == 0 {
		return ""
	}

	if len(v.Build) > 0 {
		if len(v.Divider) > 0 {
			return fmt.Sprintf("%s%s%s", v.Version, v.Divider, v.Build)
		}

		return fmt.Sprintf("%s.%s", v.Version, v.Build)
	}

	return v.Version
}

// Extract ...
func Extract(field, regex string) (Version, error) {
	tag, err := regexp.Compile(regex)
	if err != nil {
		return Version{}, err
	}

	matches := tag.FindStringSubmatch(field)

	return Version{
		Version: getPart("version", matches),
		Build:   getPart("build", matches),
		Divider: getPart("divider", matches),
	}, nil
}

func getPart(part string, matches []string) string {

	vIndex := map[string]int{
		"full":    0,
		"version": 1,
		"divider": 2,
		"build":   3,
	}

	if len(matches) == 3 {
		delete(vIndex, "divider")
		vIndex["build"] = 2
	}

	i, found := vIndex[part]

	if !found || i >= len(matches) {
		return ""
	}

	return matches[i]
}

package semver

import (
	"fmt"
	"regexp"
)

// The regular expression that represents a semantic version
var r = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type Semver struct {
	Major, Minor, Patch       string
	PreRelease, BuildMetaData *string
}

// Parse parses a string to a semantic version
func Parse(str string) (Semver, error) {
	matches := r.FindStringSubmatch(str)
	if matches == nil {
		return Semver{}, fmt.Errorf("%s is not a valid semver", str)
	}
	names := r.SubexpNames()
	result := make(map[string]string)
	for i, match := range matches {
		if i != 0 {
			result[names[i]] = match
		}
	}
	s := Semver{
		Major:         result["major"],
		Minor:         result["minor"],
		Patch:         result["patch"],
		PreRelease:    nil,
		BuildMetaData: nil,
	}

	if result["prerelease"] != "" {
		v := result["prerelease"]
		s.PreRelease = &v
	}
	if result["buildmetadata"] != "" {
		v := result["buildmetadata"]
		s.BuildMetaData = &v
	}
	return s, nil
}

func (s Semver) String() string {
	// Start with the base version
	result := fmt.Sprintf("%s.%s.%s", s.Major, s.Minor, s.Patch)

	// Append pre-release version if it exists
	if s.PreRelease != nil {
		result = fmt.Sprintf("%s-%s", result, *s.PreRelease)
	}

	// Append build metadata if it exists
	if s.BuildMetaData != nil {
		result = fmt.Sprintf("%s+%s", result, *s.BuildMetaData)
	}

	return result
}

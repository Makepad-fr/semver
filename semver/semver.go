package semver

import (
	"fmt"
	"regexp"
	"strconv"
)

// The regular expression that represents a semantic version
var r = regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type Semver struct {
	Major, Minor, Patch       string
	PreRelease, BuildMetaData *string
}

type SemverDifference struct {
	Major, Minor, Patch       int
	PreRelease, BuildMetaData []StringDifference
}

type StringDifference struct {
	Added      bool
	StartIndex uint
	Content    string
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

// String returns the string representation of the current Semver
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

// Diff calculates the difference between current Semver and another Semver
func (s Semver) Diff(nv Semver) (SemverDifference, error) {
	return Difference(s, nv)
}

// Difference calculates the difference between two SemVer structures
func Difference(old, new Semver) (SemverDifference, error) {
	oldMajor, err := strconv.Atoi(old.Major)
	if err != nil {
		return SemverDifference{}, err
	}
	newMajor, err := strconv.Atoi(new.Major)
	if err != nil {
		return SemverDifference{}, err
	}

	oldMinor, err := strconv.Atoi(old.Minor)
	if err != nil {
		return SemverDifference{}, err
	}
	newMinor, err := strconv.Atoi(new.Minor)
	if err != nil {
		return SemverDifference{}, err
	}

	oldPatch, err := strconv.Atoi(old.Patch)
	if err != nil {
		return SemverDifference{}, err
	}
	newPatch, err := strconv.Atoi(new.Patch)
	if err != nil {
		return SemverDifference{}, err
	}

	preReleaseDiff := differenceStringPtr(old.PreRelease, new.PreRelease)
	buildMetaDataDiff := differenceStringPtr(old.BuildMetaData, new.BuildMetaData)

	diff := SemverDifference{
		Major:         newMajor - oldMajor,
		Minor:         newMinor - oldMinor,
		Patch:         newPatch - oldPatch,
		PreRelease:    preReleaseDiff,
		BuildMetaData: buildMetaDataDiff,
	}
	return diff, nil
}

func (s Semver) FlatPreReleaseDifference(new Semver) (int, error) {
	if s.PreRelease != nil && new.PreRelease != nil {
		oldPR, err := strconv.Atoi(*s.PreRelease)
		if err == nil {
			newPR, err := strconv.Atoi(*new.PreRelease)
			if err == nil {
				return newPR - oldPR, nil
			}
		}
	}
	return -1, fmt.Errorf("Can not calculate flat difference between pre-release parts of %s and %s", s, new)
}

// differenceStringPtr calculates the difference of two string pointers
func differenceStringPtr(old, new *string) []StringDifference {
	if old == nil && new == nil {
		return []StringDifference{}
	}
	if old == nil {
		return []StringDifference{
			{
				StartIndex: 0,
				Added:      true,
				Content:    *new,
			},
		}
	}
	if new == nil {
		return []StringDifference{
			{
				StartIndex: 0,
				Added:      false,
				Content:    *old,
			},
		}
	}
	var result []StringDifference
	minLength := min(len(*old), len(*new))
	var addedString, removedString []byte
	startIndex := -1
	var i = 0
	for i = 0; i < minLength; i++ {
		if (*old)[i] != (*new)[i] {
			if startIndex == -1 {
				startIndex = i
				addedString = make([]byte, 0, 1)
				removedString = make([]byte, 0, 1)
			}
			addedString = append(addedString, (*new)[i])
			removedString = append(removedString, (*old)[i])
			continue
		}
		if startIndex != -1 {
			result = append(result, StringDifference{
				StartIndex: uint(startIndex),
				Added:      false,
				Content:    string(removedString),
			})
			result = append(result, StringDifference{
				StartIndex: uint(startIndex),
				Added:      true,
				Content:    string(addedString),
			})
			startIndex = -1
			addedString = nil
			removedString = nil
		}
	}
	if startIndex != -1 {
		result = append(result, StringDifference{
			StartIndex: uint(startIndex),
			Added:      false,
			Content:    string(removedString),
		})
		result = append(result, StringDifference{
			StartIndex: uint(startIndex),
			Added:      true,
			Content:    string(addedString),
		})
		startIndex = -1
		addedString = nil
		removedString = nil
	}

	if len(*old) > minLength {
		// If old string is longer
		result = append(result, StringDifference{
			StartIndex: uint(i),
			Added:      false,
			Content:    (*old)[i:],
		})
	}

	if len(*new) > minLength {
		// If new string is longer
		result = append(result, StringDifference{
			StartIndex: uint(i),
			Added:      true,
			Content:    (*new)[i:],
		})
	}
	return result
}

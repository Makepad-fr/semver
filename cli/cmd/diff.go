package cmd

import (
	"errors"
	"flag"
	"fmt"
	"strings"
)

var DiffCommand = Command{
	Name:        "diff",
	Description: "",
	Run: func() error {
		var oldVersionFlag semverFlag
		flag.Var(&oldVersionFlag, "old", "Old version")
		var lastVersionFlag semverFlag
		flag.Var(&lastVersionFlag, "new", "New version")

		var expectedPatchVersionDiff, expectedMinorVersionDiff, expectedMajorVersionDiff, expectedPreReleaseDiff intFlag

		flag.Var(&expectedPatchVersionDiff, "expected-patch-diff", "Expected patch version difference")
		flag.Var(&expectedMinorVersionDiff, "expected-minor-diff", "Expected major difference")
		flag.Var(&expectedMajorVersionDiff, "expected-major-diff", "Expected major difference")
		flag.Var(&expectedPreReleaseDiff, "expected-pre-release-diff", "Expected pre-release difference")

		var preReleasePrefixFlag, preReleaseSuffixFlag stringFlag
		flag.Var(&preReleasePrefixFlag, "pre-release-prefix", "The prefix for pre-release part of both versions")
		flag.Var(&preReleaseSuffixFlag, "pre-release-suffix", "The suffix for pre-release part of both versions")

		flag.Parse()

		if !oldVersionFlag.set {
			return errors.New("-old option for old version is required")
		}
		if !lastVersionFlag.set {
			return errors.New("-new option for new version is required")
		}

		var opr, npr string
		if preReleasePrefixFlag.set {
			if oldVersionFlag.value.PreRelease == nil {
				return fmt.Errorf("%s does not have a pre-release part but pre-release-prefix flag is passed with %s",
					oldVersionFlag.value, preReleasePrefixFlag.value)
			}
			if !strings.HasPrefix(*oldVersionFlag.value.PreRelease, preReleasePrefixFlag.value) {
				return fmt.Errorf("%s does not have a pre-release prefix %s",
					oldVersionFlag.value, preReleasePrefixFlag.value)
			}
			if lastVersionFlag.value.PreRelease == nil {
				return fmt.Errorf("%s does not have a pre-release part but pre-release-prefix flag is passed with %s",
					lastVersionFlag.value, preReleasePrefixFlag.value)
			}
			if !strings.HasPrefix(*lastVersionFlag.value.PreRelease, preReleasePrefixFlag.value) {
				return fmt.Errorf("%s does not have a pre-release prefix %s",
					lastVersionFlag.value, preReleasePrefixFlag.value)
			}
			opr = strings.TrimPrefix(*oldVersionFlag.value.PreRelease, preReleasePrefixFlag.value)
			npr = strings.TrimPrefix(*lastVersionFlag.value.PreRelease, preReleasePrefixFlag.value)
		}
		if preReleaseSuffixFlag.set {
			if oldVersionFlag.value.PreRelease == nil {
				return fmt.Errorf("%s does not have a pre-release part but pre-release-suffix flag is passed with %s",
					oldVersionFlag.value, preReleaseSuffixFlag.value)
			}
			if !strings.HasSuffix(*oldVersionFlag.value.PreRelease, preReleaseSuffixFlag.value) {
				return fmt.Errorf("%s does not have a pre-release prefix %s",
					oldVersionFlag.value, preReleaseSuffixFlag.value)
			}
			if lastVersionFlag.value.PreRelease == nil {
				return fmt.Errorf("%s does not have a pre-release part but pre-release-prefix flag is passed with %s",
					lastVersionFlag.value, preReleaseSuffixFlag.value)
			}
			if !strings.HasPrefix(*lastVersionFlag.value.PreRelease, preReleaseSuffixFlag.value) {
				return fmt.Errorf("%s does not have a pre-release prefix %s",
					&lastVersionFlag.value, preReleaseSuffixFlag.value)
			}
			if len(strings.TrimSpace(opr)) == 0 {
				opr = *oldVersionFlag.value.PreRelease
				npr = *lastVersionFlag.value.PreRelease
			}
			opr = strings.TrimSuffix(opr, preReleaseSuffixFlag.value)
			npr = strings.TrimPrefix(npr, preReleaseSuffixFlag.value)
		}

		diff, err := oldVersionFlag.value.Diff(lastVersionFlag.value)
		if err != nil {
			fmt.Printf("Error while calculating diff between %s and %s", oldVersionFlag.value, lastVersionFlag.value)
			return err
		}
		if expectedMajorVersionDiff.set && expectedMajorVersionDiff.value != diff.Major {
			return fmt.Errorf("expected major version dfference is %d: calculated major version difference is %d", expectedMajorVersionDiff.value, diff.Major)
		}
		if expectedMinorVersionDiff.set && expectedMinorVersionDiff.value != diff.Minor {
			return fmt.Errorf("expected minor version dfference is %d: calculated minor version difference is %d", expectedMinorVersionDiff.value, diff.Minor)
		}
		if expectedPatchVersionDiff.set && expectedPatchVersionDiff.value != diff.Patch {
			return fmt.Errorf("expected patch version dfference is %d: calculated patch version difference is %d", expectedPatchVersionDiff.value, diff.Patch)
		}

		if expectedPreReleaseDiff.set {
			oprc := oldVersionFlag.value.PreRelease
			oldVersionFlag.value.PreRelease = &opr
			nprc := lastVersionFlag.value.PreRelease
			lastVersionFlag.value.PreRelease = &npr
			diff, err := oldVersionFlag.value.FlatPreReleaseDifference(lastVersionFlag.value)
			if err != nil {
				return err
			}
			oldVersionFlag.value.PreRelease = oprc
			lastVersionFlag.value.PreRelease = nprc

			if diff != expectedPreReleaseDiff.value {
				return fmt.Errorf("pre-release difference between %s and %s is %d instead of %d", &oldVersionFlag.value, lastVersionFlag.value, diff, expectedPreReleaseDiff.value)
			}
		}

		return nil
	},
}

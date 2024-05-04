package semver

import (
	"fmt"
	"testing"
)

var validSemvers = []string{
	"0.0.4",
	"1.2.3",
	"10.20.30",
	"1.1.2-prerelease+meta",
	"1.1.2+meta",
	"1.1.2+meta-valid",
	"1.0.0-alpha",
	"1.0.0-beta",
	"1.0.0-alpha.beta",
	"1.0.0-alpha.beta.1",
	"1.0.0-alpha.1",
	"1.0.0-alpha0.valid",
	"1.0.0-alpha.0valid",
	"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay",
	"1.0.0-rc.1+build.1",
	"2.0.0-rc.1+build.123",
	"1.2.3-beta",
	"10.2.3-DEV-SNAPSHOT",
	"1.2.3-SNAPSHOT-123",
	"1.0.0",
	"2.0.0",
	"1.1.7",
	"2.0.0+build.1848",
	"2.0.1-alpha.1227",
	"1.0.0-alpha+beta",
	"1.2.3----RC-SNAPSHOT.12.9.1--.12+788",
	"1.2.3----R-S.12.9.1--.12+meta",
	"1.2.3----RC-SNAPSHOT.12.9.1--.12",
	"1.0.0+0.build.1-rc.10000aaa-kk-0.1",
	"99999999999999999999999.999999999999999999.99999999999999999",
	"1.0.0-0A.is.legal",
}

var invalidSemvers = []string{
	"1",
	"1.2",
	"1.2.3-0123",
	"1.2.3-0123.0123",
	"1.1.2+.123",
	"+invalid",
	"-invalid",
	"-invalid+invalid",
	"-invalid.01",
	"alpha",
	"alpha.beta",
	"alpha.beta.1",
	"alpha.1",
	"alpha+beta",
	"alpha_beta",
	"alpha.",
	"alpha..",
	"beta",
	"1.0.0-alpha_beta",
	"-alpha.",
	"1.0.0-alpha..",
	"1.0.0-alpha..1",
	"1.0.0-alpha...1",
	"1.0.0-alpha....1",
	"1.0.0-alpha.....1",
	"1.0.0-alpha......1",
	"1.0.0-alpha.......1",
	"01.1.1",
	"1.01.1",
	"1.1.01",
	"1.2",
	"1.2.3.DEV",
	"1.2-SNAPSHOT",
	"1.2.31.2.3----RC-SNAPSHOT.12.09.1--..12+788",
	"1.2-RC-SNAPSHOT",
	"-1.0.3-gamma+b7718",
	"+justmeta",
	"9.8.7+meta+meta",
	"9.8.7-whatever+meta+meta",
	"99999999999999999999999.999999999999999999.99999999999999999----RC-SNAPSHOT.12.09.1--------------------------------..12",
}

func TestParseValidSemvers(t *testing.T) {
	t.Parallel()
	for _, v := range validSemvers {
		t.Run(fmt.Sprintf("%s should parse", v), func(t *testing.T) {
			t.Parallel()
			_, err := Parse(v)
			if err != nil {
				t.Errorf("Error while parsing valid semver %s: %v", v, err)
				return
			}
		})
	}
}

func TestParseInvalidSemvers(t *testing.T) {
	t.Parallel()
	for _, v := range invalidSemvers {
		t.Run(fmt.Sprintf("%s should not be parsed", v), func(t *testing.T) {
			t.Parallel()
			_, err := Parse(v)
			if err == nil {
				t.Errorf("Parsing an invalid semver %s should fail", v)
				return
			}
		})
	}
}

func TestDifferenceBaseVersion(t *testing.T) {
	t.Parallel()
	old, _ := Parse("1.0.0")
	new, _ := Parse("2.0.0")
	diff, err := old.Diff(new)
	if err != nil {
		t.Errorf("Error while calculating difference between %s and %s: %v", old, new, err)
	}
	if diff.Major != 1 {
		t.Errorf("Major difference should 1: %d", diff.Major)
	}
	if diff.Minor != 0 {
		t.Errorf("Minor difference should 0: %d", diff.Minor)
	}
	if diff.Patch != 0 {
		t.Errorf("Patch difference should 0: %d", diff.Patch)
	}
	if len(diff.PreRelease) != 0 {
		t.Errorf("PreRelease difference should be empty")
	}
	if len(diff.BuildMetaData) != 0 {
		t.Errorf("BuildMetedata difference should be empty")
	}
}
func TestDifferencePreRelease(t *testing.T) {
	t.Parallel()
	old, _ := Parse("1.0.0-rc1")
	new, _ := Parse("1.0.0-rc2")
	diff, err := old.Diff(new)
	if err != nil {
		t.Errorf("Error while calculating difference between %s and %s: %v", old, new, err)
	}
	if diff.Major != 0 {
		t.Errorf("Major difference should 0: %d", diff.Major)
	}
	if diff.Minor != 0 {
		t.Errorf("Minor difference should 0: %d", diff.Minor)
	}
	if diff.Patch != 0 {
		t.Errorf("Patch difference should 0: %d", diff.Patch)
	}
	if len(diff.PreRelease) != 2 {
		t.Errorf("PreRelease difference should have 2 elements")

	}
	if len(diff.BuildMetaData) != 0 {
		t.Errorf("BuildMetedata difference should be empty")
	}
	if !(diff.PreRelease[0].Content == "1" && !diff.PreRelease[0].Added && diff.PreRelease[0].StartIndex == 2) {
		t.Errorf("First pre-release difference content should be 1: %s, it should be added false = %v, its index should be 8: %d",
			diff.PreRelease[0].Content,
			diff.PreRelease[0].Added,
			diff.PreRelease[0].StartIndex,
		)
	}
	if !(diff.PreRelease[1].Content == "2" && diff.PreRelease[1].Added && diff.PreRelease[1].StartIndex == 2) {
		t.Errorf("Second pre-release difference content should be 2: %s, it should be added true = %v, its index should be 8: %d",
			diff.PreRelease[1].Content,
			diff.PreRelease[1].Added,
			diff.PreRelease[1].StartIndex,
		)
	}
	t.Logf("Second pre-release difference %s", diff.PreRelease[1].Content)
}

package cmd

import (
	"github.com/Makepad-fr/semver/semver"
)

type semverFlag struct {
	set   bool
	value semver.Semver
}

func (i *semverFlag) String() string {
	return i.value.String()
}

func (i *semverFlag) Set(s string) error {
	v, err := semver.Parse(s)
	if err != nil {
		return err
	}
	i.set = true
	i.value = v
	return nil
}

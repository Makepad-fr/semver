package cmd

import "strings"

type stringFlag struct {
	set   bool
	value string
}

func (i *stringFlag) String() string {
	return i.value
}

func (i *stringFlag) Set(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}
	i.set = true
	i.value = s
	return nil
}

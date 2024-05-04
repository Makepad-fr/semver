package cmd

import (
	"fmt"
	"strconv"
)

type intFlag struct {
	set   bool
	value int
}

func (i *intFlag) String() string {
	return fmt.Sprintf("%s", i.value)
}

func (i *intFlag) Set(s string) error {
	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	i.set = true
	i.value = v
	return nil
}

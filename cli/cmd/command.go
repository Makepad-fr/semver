package cmd

import "fmt"

type Command struct {
	Name        string
	Description string
	Run         func() error
}

var commands map[string]Command

func init() {
	commands = make(map[string]Command)
}

func AddCommand(c Command) error {
	if _, ok := commands[c.Name]; ok {
		return fmt.Errorf("Command %s already added", c.Name)
	}
	commands[c.Name] = c
	return nil
}

func GetCommand(name string) (Command, bool) {
	c, ok := commands[name]
	if ok {
		return c, true
	}
	return Command{}, false
}

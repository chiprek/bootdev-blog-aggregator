package main

import (
	"context"
	"fmt"
)

type commands struct {
	cmds map[string]func(*state, command) error
}

type command struct {
	name string
	args []string
}

func (c *commands) run(s *state, cmd command) error {
	toRun, ok := c.cmds[cmd.name]
	if !ok {
		return fmt.Errorf("No such command: %s\n", cmd.name)
	}
	return toRun(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) error {
	_, ok := c.cmds[name]
	if ok {
		return fmt.Errorf("Command %s already exists\n", name)
	}
	c.cmds[name] = f
	return nil
}

func handlerReset(s *state, cmd command) error {
	var err error

	err = s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}

	return err
}

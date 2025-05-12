package main

import (
	"errors"
	"fmt"

	"github.com/janiv/gator/internal/config"
)

type State struct {
	cfg *config.Config
}

func NewState() *State {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	return &State{
		cfg: &c,
	}
}

type Command struct {
	name string
	args []string
}

type Commands struct {
	commands map[string]func(*State, Command) error
}

func (c *Commands) run(s *State, cmd Command) error {
	val, exists := c.commands[cmd.name]
	if !exists {
		return errors.New("yo you missing a command")
	}
	val(s, cmd)
	return nil
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	err := s.cfg.SetUser(cmd.args[1])
	if err != nil {
		return err
	}
	fmt.Println("user set")
	return nil
}

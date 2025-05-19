package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/janiv/gator/internal/config"
	"github.com/janiv/gator/internal/database"
)

type State struct {
	db  *database.Queries
	cfg *config.Config
}

func NewState() *State {
	c, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return nil
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
	return val(s, cmd)
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	c.commands[name] = f
}

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	name := cmd.args[0]
	usr_check, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("%s is not in database\n", name)
		return err
	}
	set_err := s.cfg.SetUser(usr_check.Name)
	if set_err != nil {
		return set_err
	}
	fmt.Println("user set")
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	fmt.Printf("args= %s\n", cmd.args[0])
	name := cmd.args[0]
	usr_check, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		fmt.Printf("%s not in database, attempting to create\n", name)
	} else {
		fmt.Printf("%s already exists\n", usr_check.Name)
		return errors.New("user already exists")
	}
	curr_time := time.Now()
	db_params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: curr_time,
		UpdatedAt: curr_time,
		Name:      name,
	}
	fmt.Println(db_params)
	usr, err := s.db.CreateUser(context.Background(), db_params)
	if err != nil {
		return err
	}
	handlerLogin(s, cmd)
	fmt.Printf("User %s was created\n", usr.Name)
	return nil
}

func handlerReset(s *State, cmd Command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("database reset")
	return nil
}

func handlerUsers(s *State, cmd Command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	// ok this is a bit stupid but we doing it anyway
	curr := users[0]

	for _, usr := range users {
		if usr.UpdatedAt.After(curr.UpdatedAt) {
			curr = usr
		}
	}
	for _, usr := range users {
		fmt.Printf("* %s ", usr.Name)
		if usr == curr {
			fmt.Print("(current)")
		}
		fmt.Print("\n")
	}
	return nil
}

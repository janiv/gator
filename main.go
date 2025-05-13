package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	s := NewState()
	com_map := make(map[string]func(*State, Command) error)
	coms := Commands{
		commands: com_map,
	}
	coms.register("login", handlerLogin)
	stuff := os.Args
	fmt.Println(len(stuff))
	if len(stuff) == 1 {
		e := errors.New("not enough arguments")
		fmt.Printf("%s\n", e)
		os.Exit(1)
	} else if len(stuff) == 2 {
		e := errors.New("username required for login")
		fmt.Printf("%s\n", e)
		os.Exit(1)
	} else if len(stuff) == 3 {
		fmt.Println(stuff)
		var args_temp []string = make([]string, 1)
		args_temp[0] = stuff[2]
		fmt.Println(args_temp)
		issue := Command{
			name: stuff[1],
			args: args_temp,
		}
		e := coms.run(s, issue)
		if e != nil {
			fmt.Printf("%s\n", e)
		}
	} else {
		e := errors.New("too much stuff")
		fmt.Printf("%s\n", e)
		os.Exit(1)
	}
}

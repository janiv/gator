package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/janiv/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	s := NewState()
	db, err := sql.Open("postgres", s.cfg.DbURL)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)
	s.db = dbQueries
	com_map := make(map[string]func(*State, Command) error)
	coms := Commands{
		commands: com_map,
	}
	coms.register("login", handlerLogin)
	coms.register("register", handlerRegister)
	stuff := os.Args
	stuff_len := len(stuff)
	switch stuff_len {
	case 1:
		{
			fmt.Println("missing arguments")
			os.Exit(1)
		}
	case 2:
		{
			fmt.Println("missing username for login/register")
			os.Exit(1)
		}
	case 3:
		{
			var args_temp []string = make([]string, 1)
			args_temp[0] = stuff[2]
			issue := Command{
				name: stuff[1],
				args: args_temp,
			}
			e := coms.run(s, issue)
			if e != nil {
				fmt.Printf("%s\n", e)
				os.Exit(1)
			} else {
				fmt.Println("Hey it worked")
				os.Exit(0)
			}
		}
	default:
		{
			fmt.Println("I don't handle essays")
			os.Exit(1)
		}
	}
}

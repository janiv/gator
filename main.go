package main

import (
	"fmt"

	"github.com/janiv/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	cfg.SetUser("jani")
	cfg2, err := config.Read()
	fmt.Println(cfg2.CurrentUserName)
	fmt.Println(cfg2.DbURL)

}

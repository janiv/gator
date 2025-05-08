package main


type state struct {
	cfg *Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("missing args")
	}
	err := cfg.SetUserName(cmd.args[1])
	if err != nil {
		return err
	}
	fmt.Println("user set")
	return nil
}


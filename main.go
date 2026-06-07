package main

import (
	"fmt"
	"os"

	"github.com/chiprek/bootdev-blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}

	s := state{cfg: &cfg}
	cmds := commands{cmds: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Command not specified")
		os.Exit(1)
	}

	cmdName := os.Args[1]
	err = cmds.run(&s, command{cmdName, os.Args[2:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

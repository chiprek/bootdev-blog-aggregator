package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/chiprek/bootdev-blog-aggregator/internal/config"
	"github.com/chiprek/bootdev-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		fmt.Printf("Failed to open db connection: %v\n", err)
		os.Exit(1)
	}

	s := state{cfg: &cfg, db: database.New(db)}
	cmds := commands{cmds: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)

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

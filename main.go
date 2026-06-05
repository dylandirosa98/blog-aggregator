package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/dylandirosa98/blog-aggregator/internal/config"
	"github.com/dylandirosa98/blog-aggregator/internal/database"
)
import _ "github.com/lib/pq"

func main() {
	con, err := config.Read()
	if err != nil {
		print(fmt.Errorf("error reading config file: %v", err))
	}
	db, err := sql.Open("postgres", con.Db_url)
	dbQueries := database.New(db)
	theState := state{
		config: &con,
		db:     dbQueries,
	}
	theCommands := &commands{mapCommands: make(map[string]func(*state, command) error)}
	theCommands.register("login", handlerLogin)
	theCommands.register("register", handlerRegister)
	theCommands.register("reset", handlerResetUser)
	theCommands.register("users", handlerGetUsers)
	theCommands.register("agg", handlerAgg)
	theCommands.register("addfeed", handlerAddFeed)
	theCommands.register("feeds", handlerFeeds)
	if err != nil {
		print(fmt.Errorf("error registering command: %v", err))
	}
	if len(os.Args) < 2 {
		if len(os.Args) > 1 {
			fmt.Printf("not enough arguments\n")
		} else {
			fmt.Printf("a username is required\n")
		}
		os.Exit(1)
	}
	cmd := command{os.Args[1], os.Args[2:]}
	err = theCommands.run(&theState, cmd)
	if err != nil {
		print(fmt.Errorf("error running command: %v", err))
		os.Exit(1)
	}
	os.Exit(0)
}

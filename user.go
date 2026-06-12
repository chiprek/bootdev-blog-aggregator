package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chiprek/bootdev-blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Login handler expects a single argument, the username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Logged in as: %s\n", cmd.args[0])
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Register handler expectes a single argument, the username")
	}

	name := cmd.args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: name})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return fmt.Errorf("User %s already exists", name)
			}
		}
		return fmt.Errorf("Failed to create user: %v", err)
	}

	s.cfg.SetUser(user.Name)
	fmt.Printf("Registered user: %s\n", user.Name)
	return nil
}

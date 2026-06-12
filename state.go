package main

import (
	"github.com/chiprek/bootdev-blog-aggregator/internal/config"
	"github.com/chiprek/bootdev-blog-aggregator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

package main

import (
	"database/sql"
	"flag"

	"github.com/pressly/goose"
	"github.com/yeralin-munar/tt-go-json-fernet/config"
	"github.com/yeralin-munar/tt-go-json-fernet/internal/data/postgres"

	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// flagconf is the config flag.
	flagconf string
)

func main() {
	flag.StringVar(&flagconf, "conf", "../configs", "config path, eg: -conf config.yaml")
	flag.Parse()
	config.NewConfig(flagconf)

	cfg := config.Cfg

	dbURL := postgres.GenerateDBURL(&cfg.Data)
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// setup database
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}
}

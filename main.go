package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/mbaev/simplebank/api"
	db "github.com/mbaev/simplebank/db/sqlc"
	"github.com/mbaev/simplebank/util"
)

func main() {
	conf, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config file", err)
	}
	conn, err := sql.Open(conf.DBDriver, conf.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(conf.ServerAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}

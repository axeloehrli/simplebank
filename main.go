package main

import (
	"database/sql"
	"log"

	"github.com/axeloehrli/simplebank/api"
	db "github.com/axeloehrli/simplebank/db/sqlc"
	"github.com/axeloehrli/simplebank/db/util"

	_ "github.com/lib/pq"
)

func main() {

	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("Cannot load config file: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Could not start server: ", err)
	}
}

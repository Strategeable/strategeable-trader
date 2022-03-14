package main

import (
	"github.com/Strategeable/Trader/database"
	"github.com/Strategeable/Trader/rpcserver"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("error loading .env file")
	}

	databaseHandler := &database.DatabaseHandler{}
	err = databaseHandler.Connect()
	if err != nil {
		panic(err)
	}

	server := rpcserver.NewRpcServer(databaseHandler)
	server.Start()

	keepaliveCh := make(chan string)
	<-keepaliveCh
}

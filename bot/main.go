package main

import (
	"github.com/Strategeable/Trader/communication"
	"github.com/Strategeable/Trader/database"
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

	communication.SetupAmqp(databaseHandler)

	keepaliveCh := make(chan string)
	<-keepaliveCh
}

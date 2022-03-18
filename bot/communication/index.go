package communication

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Strategeable/Trader/database"
	"github.com/Strategeable/Trader/handlers"
	"github.com/Strategeable/Trader/types"
	"github.com/streadway/amqp"
)

func SetupAmqp(databaseHandler *database.DatabaseHandler) {
	connection, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		panic(err)
	}

	go startBacktestConsumer(connection, databaseHandler)
}

func startBacktestConsumer(connection *amqp.Connection, databaseHandler *database.DatabaseHandler) {
	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	// Handle 1 backtest at once
	channel.Qos(1, 0, false)

	// Ensure required exchanges exist
	err = channel.ExchangeDeclare("backtest_x", "topic", false, false, false, false, amqp.Table{})
	if err != nil {
		panic(err)
	}

	// Ensure queues exist
	backtestQueue, err := channel.QueueDeclare("backtest_q", true, false, false, false, amqp.Table{})
	if err != nil {
		panic(err)
	}

	backtestCh, err := channel.Consume(backtestQueue.Name, "", false, false, false, false, amqp.Table{})
	if err != nil {
		panic(err)
	}

	backtestHandler := handlers.NewBacktestHandler(databaseHandler)

	for delivery := range backtestCh {
		backtestId := &types.QueuedBacktest{}
		err := json.Unmarshal(delivery.Body, backtestId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ch, err := backtestHandler.RunBacktest(backtestId.Id)
		if err != nil {
			fmt.Println(err)
			channel.Ack(delivery.DeliveryTag, false)
			continue
		}

		backtestKey := fmt.Sprintf("backtests.%s", backtestId.Id)

		fmt.Printf("Running backtest %s. Sending updates to %s.\n", backtestId.Id, backtestKey)

		status := ""

		for event := range ch {
			body, err := json.Marshal(event)
			if err != nil {
				fmt.Println(err)
				continue
			}

			channel.Publish("backtest_x", backtestKey, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})

			if event.Status != status {
				fmt.Printf("Backtest %s - Status: %s.\n", backtestId.Id, event.Status)
				status = event.Status
				databaseHandler.UpdateBacktestStatus(backtestId.Id, event.Status)
			}
		}

		fmt.Printf("Finished backtest %s.\n", backtestId.Id)

		channel.Ack(delivery.DeliveryTag, false)
	}
}
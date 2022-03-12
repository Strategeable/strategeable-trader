package rpcserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"

	"github.com/Stratomicl/Trader/database"
	"github.com/Stratomicl/Trader/handlers"
	"github.com/Stratomicl/Trader/impl"
	strategy_types "github.com/Stratomicl/Trader/strategy"
	"github.com/Stratomicl/Trader/types"
)

type rpcServer struct {
	server          *rpc.Server
	databaseHandler *database.DatabaseHandler
}

func NewRpcServer(databaseHandler *database.DatabaseHandler) *rpcServer {
	return &rpcServer{
		server:          rpc.NewServer(),
		databaseHandler: databaseHandler,
	}
}

func (r *rpcServer) Start() error {
	backtest := &Backtest{
		databaseHandler: r.databaseHandler,
	}

	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(backtest, "Backtest")
	http.Handle("/rpc", s)
	http.ListenAndServe("localhost:1234", nil)

	return nil
}

type Backtest struct {
	databaseHandler *database.DatabaseHandler
}

func (b *Backtest) Backtest(r *http.Request, backtestId *string, reply *int) error {
	backtest, err := b.databaseHandler.GetBacktestById(*backtestId)
	if err != nil {
		return err
	}
	*reply = 1

	go b.performBacktest(backtest)
	return nil
}

func (b *Backtest) performBacktest(backtest *strategy_types.Backtest) {
	strategy, err := strategy_types.StrategyFromJson(backtest.Strategy)
	if err != nil {
		panic(err)
	}

	exchangeImpl := impl.NewBinanceExchangeImpl()

	err = exchangeImpl.Init()
	if err != nil {
		panic(err)
	}

	from := backtest.FromDate.Time()
	to := backtest.ToDate.Time()

	marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, strategy.Symbols, strategy.GetTimeFrames())

	positionHandler := impl.NewSimulatedPositionHandler(backtest.StartBalance, make([]*types.Position, 0))

	eventCh := make(chan types.PositionHandlerEvent, 5)
	positionHandler.SubscribeEvents(eventCh)

	engine := handlers.NewEngine(strategy, marketDataProvider, positionHandler)

	finishCh := make(chan string)

	go func() {
		err = engine.Start()
		if err != nil {
			panic(err)
		}

		close(finishCh)
	}()

	positions := make([]*types.Position, 0)

	for {
		select {
		case event := <-eventCh:
			switch event.Type {
			case types.POSITION_CREATED:
				position := event.Data.(*types.Position)

				positions = append(positions, position)

				fmt.Printf("[BACKTEST] Position created: %s at %.2f.\n", position.Symbol().String(), position.AverageEntryRate())
			case types.POSITION_CLOSED:
				position := event.Data.(*types.Position)
				fmt.Printf("[BACKTEST] Position closed: %s at %.2f. Change %%: %.2f.\n", position.Symbol().String(), position.AverageExitRate(0), position.ChangePercentage(0))
			}
		case _, ok := <-finishCh:
			if ok {
				continue
			}

			mappedPositions := make([]strategy_types.BacktestPosition, 0)
			for _, position := range positions {
				closeTime := time.Unix(0, 0)
				if position.IsClosed() {
					closeTime = *position.CloseTime()
				}
				mappedPositions = append(mappedPositions, strategy_types.BacktestPosition{
					OpenedAt: position.OpenTime(),
					ClosedAt: closeTime,
					Symbol:   position.Symbol().String(),
					EntryValue: strategy_types.BacktestPositionValue{
						Date:      position.OpenTime(),
						Rate:      position.AverageEntryRate(),
						BaseSize:  position.BaseSize(),
						QuoteFees: position.EntryQuoteFees(),
					},
					ExitValue: strategy_types.BacktestPositionValue{
						Date:      closeTime,
						Rate:      position.AverageExitRate(0),
						BaseSize:  position.BaseSize(),
						QuoteFees: position.ExitQuoteFees(),
					},
				})
			}

			backtest.Finished = true
			backtest.EndBalance = positionHandler.TotalBalance
			backtest.Positions = mappedPositions

			err := b.databaseHandler.SaveBacktest(backtest)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%.2f => %.2f\n", backtest.StartBalance, positionHandler.TotalBalance)
			return
		}
	}
}

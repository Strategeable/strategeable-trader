package rpcserver

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/Stratomicl/Trader/database"
	"github.com/Stratomicl/Trader/handlers"
	"github.com/Stratomicl/Trader/impl"
	"github.com/Stratomicl/Trader/strategy"
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

	rpc.Register(backtest)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		return e
	}
	go http.Serve(l, nil)

	return nil
}

type Backtest struct {
	databaseHandler *database.DatabaseHandler
}

func (b *Backtest) Backtest(backtestId string, reply *int) error {
	backtest, err := b.databaseHandler.GetBacktestById(backtestId)
	if err != nil {
		return err
	}
	*reply = 1

	go b.performBacktest(backtest)
	return nil
}

func (b *Backtest) performBacktest(backtest *strategy.Backtest) {
	strategy, err := strategy.StrategyFromJson(backtest.Strategy)
	if err != nil {
		panic(err)
	}

	exchangeImpl := impl.NewBinanceExchangeImpl()

	err = exchangeImpl.Init()
	if err != nil {
		panic(err)
	}

	from := backtest.DateFrom
	to := backtest.DateTo

	marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, strategy.Symbols, strategy.GetTimeFrames())

	positionHandler := impl.NewSimulatedPositionHandler(backtest.StartBalance, make([]*types.Position, 0))

	eventCh := make(chan types.PositionHandlerEvent, 5)
	positionHandler.SubscribeEvents(eventCh)

	engine := handlers.NewEngine(strategy, marketDataProvider, positionHandler)

	err = engine.Start()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f => %.2f\n", backtest.StartBalance, positionHandler.TotalBalance)
}

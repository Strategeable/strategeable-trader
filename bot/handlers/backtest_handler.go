package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/Strategeable/Trader/database"
	"github.com/Strategeable/Trader/impl"
	strategy_types "github.com/Strategeable/Trader/strategy"
	"github.com/Strategeable/Trader/types"
)

type BacktestHandler struct {
	databaseHandler      *database.DatabaseHandler
	mainCandleCollection *types.CandleCollection
}

func NewBacktestHandler(databaseHandler *database.DatabaseHandler) *BacktestHandler {
	return &BacktestHandler{
		databaseHandler:      databaseHandler,
		mainCandleCollection: types.NewCandleCollection(-1),
	}
}

func (b *BacktestHandler) RunBacktest(id string) (chan types.BacktestEvent, error) {
	backtest, err := b.databaseHandler.GetBacktestById(id)
	if err != nil {
		return nil, err
	}

	if backtest.Finished {
		return nil, errors.New("already finished")
	}

	backtestEventCh := make(chan types.BacktestEvent)

	go func() {
		defer close(backtestEventCh)

		strategy, err := strategy_types.StrategyFromJson(backtest.Strategy)
		if err != nil {
			return
		}

		backtestEventCh <- types.BacktestEvent{
			Status: "PREPARING_PARTS",
		}

		exchangeImpl := impl.NewBinanceExchangeImpl()

		err = exchangeImpl.Init()
		if err != nil {
			return
		}

		from := backtest.FromDate.Time()
		to := backtest.ToDate.Time()

		marketDataProvider := impl.NewHistoricalMarketDataProvider(exchangeImpl, from, to, strategy.Symbols, strategy.GetTimeFrames(), b.mainCandleCollection, b.databaseHandler)

		positionHandler := impl.NewSimulatedPositionHandler(backtest.StartBalance, make([]*types.Position, 0))

		eventCh := make(chan types.PositionHandlerEvent, 25)
		positionHandler.SubscribeEvents(eventCh)

		engine := NewEngine(strategy, marketDataProvider, positionHandler)

		go func() {
			defer close(eventCh)

			backtestEventCh <- types.BacktestEvent{
				Status: "LOADING_CANDLES",
			}
			err = engine.InitializeMarketData()
			if err != nil {
				fmt.Println(err)
				return
			}

			backtestEventCh <- types.BacktestEvent{
				Status: "RUNNING",
			}
			err = engine.Start()
			if err != nil {
				fmt.Println(err)
			}

			backtestEventCh <- types.BacktestEvent{
				Status: "FINISHED",
			}
		}()

		progressCh := make(chan types.PositionHandlerEvent)

		go func() {
			for event := range eventCh {
				progressCh <- event

				eventData := &types.BacktestEventData{
					Type: event.Type,
					Data: event.Data,
				}

				if event.Type == types.POSITION_CREATED || event.Type == types.POSITION_CLOSED {
					position := event.Data.(*types.Position)
					eventData.Data = positionToBacktestPosition(position)
				}

				backtestEventCh <- types.BacktestEvent{
					Status:    "RUNNING",
					EventData: eventData,
				}
			}
			close(progressCh)
		}()

		b.handleBacktestProgress(backtest, positionHandler, progressCh)
	}()

	return backtestEventCh, nil
}

func (b *BacktestHandler) handleBacktestProgress(backtest *strategy_types.Backtest, positionHandler types.PositionHandler, eventCh chan types.PositionHandlerEvent) {
	positions := make([]*types.Position, 0)

	for event := range eventCh {
		if event.Type == types.POSITION_CREATED {
			position := event.Data.(*types.Position)
			positions = append(positions, position)
		}
	}

	mappedPositions := make([]strategy_types.BacktestPosition, 0)
	for _, position := range positions {
		mappedPositions = append(mappedPositions, positionToBacktestPosition(position))
	}

	backtest.Finished = true
	backtest.EndBalance = positionHandler.GetTotalBalance()
	backtest.Positions = mappedPositions

	err := b.databaseHandler.SaveBacktest(backtest)
	if err != nil {
		fmt.Println(err)
	}
}

func positionToBacktestPosition(position *types.Position) strategy_types.BacktestPosition {
	closeTime := time.Unix(0, 0)
	if position.IsClosed() {
		closeTime = *position.CloseTime()
	}
	return strategy_types.BacktestPosition{
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
	}
}

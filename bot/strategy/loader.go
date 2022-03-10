package strategy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Stratomicl/Trader/helpers"
	"github.com/Stratomicl/Trader/types"
)

type rawIndicator struct {
	TimeFrame    string
	CandleBack   int
	RealTime     bool
	Offset       float64
	IndicatorKey string
	Data         map[string]interface{}
}

type rawSignalTile struct {
	Id          string
	Name        string
	Operand     Operand
	Persistence int
	IndicatorA  rawIndicator
	IndicatorB  rawIndicator
}

type rawAnySignalTile struct {
	Amount  int
	Signals []rawSignalTile
}

type rawStep struct {
	Id   string
	Type string
	Data interface{}
}

type rawChunk struct {
	Id    string
	Name  string
	Steps []rawStep
}

type rawPath struct {
	Id    string
	Name  string
	Type  string
	Steps []rawStep
}

type rawStrategy struct {
	Name    string
	Symbols []string
	Chunks  []rawChunk
	Paths   []rawPath
}

func StrategyFromJson(input string) (*Strategy, error) {
	var strategy rawStrategy
	json.Unmarshal([]byte(input), &strategy)

	chunkMapping := make(map[string][]Tile)
	for _, chunk := range strategy.Chunks {
		tiles, err := stepsToTiles(chunk.Steps, chunkMapping)
		if err != nil {
			return nil, err
		}
		chunkMapping[chunk.Id] = tiles
	}

	buyPaths := make([]*Path, 0)
	sellPaths := make([]*Path, 0)

	for _, rawPath := range strategy.Paths {
		tiles, err := stepsToTiles(rawPath.Steps, chunkMapping)
		if err != nil {
			return nil, err
		}

		path := &Path{
			Tiles: tiles,
		}

		if rawPath.Type == "BUY" {
			buyPaths = append(buyPaths, path)
		} else {
			sellPaths = append(sellPaths, path)
		}
	}

	symbols := make([]types.Symbol, 0)
	for _, symbol := range strategy.Symbols {
		split := strings.Split(symbol, "/")
		symbols = append(symbols, types.Symbol{
			BaseAsset:  split[0],
			QuoteAsset: split[1],
		})
	}

	return &Strategy{
		Exchange:         types.BINANCE,
		BuyPaths:         buyPaths,
		SellPaths:        sellPaths,
		Symbols:          symbols,
		DefaultTimeFrame: types.H1,
		BuyCooldown:      60 * time.Second,
		BuySize:          100 / float64(len(symbols)),
	}, nil
}

func stepsToTiles(steps []rawStep, chunkMapping map[string][]Tile) ([]Tile, error) {
	tiles := make([]Tile, 0)
	for _, rawStep := range steps {
		addedTiles, err := stepToTile(rawStep, chunkMapping)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, addedTiles...)
	}
	return tiles, nil
}

func stepToTile(step rawStep, chunkMapping map[string][]Tile) ([]Tile, error) {
	switch step.Type {
	case "SIGNAL_TILE":
		var data rawSignalTile

		marshalled, _ := json.Marshal(step.Data)
		json.Unmarshal(marshalled, &data)

		signalTile, err := rawSignalTileToSignalTile(data)
		if err != nil {
			return nil, err
		}

		return []Tile{
			signalTile,
		}, nil
	case "ANY_SIGNAL_TILE":
		var data rawAnySignalTile

		marshalled, _ := json.Marshal(step.Data)
		json.Unmarshal(marshalled, &data)

		signalTiles := make([]*SignalTile, 0)
		for _, rawSignalTile := range data.Signals {
			signalTile, err := rawSignalTileToSignalTile(rawSignalTile)
			if err != nil {
				return make([]Tile, 0), nil
			}
			signalTiles = append(signalTiles, signalTile)
		}

		return []Tile{
			&AnySignalTile{
				Amount:      data.Amount,
				SignalTiles: signalTiles,
			},
		}, nil
	case "CHUNK_ID":
		return chunkMapping[step.Data.(string)], nil
	}

	return nil, fmt.Errorf("tile type not found: %s", step.Type)
}

func rawSignalTileToSignalTile(raw rawSignalTile) (*SignalTile, error) {
	indicatorA, err := rawIndicatorToIndicator(raw.IndicatorA)
	if err != nil {
		return nil, err
	}
	indicatorB, err := rawIndicatorToIndicator(raw.IndicatorB)
	if err != nil {
		return nil, err
	}

	timeFrameA := types.TimeFrame(raw.IndicatorA.TimeFrame)
	if len(timeFrameA) == 0 {
		timeFrameA = types.M1
	}
	timeFrameB := types.TimeFrame(raw.IndicatorB.TimeFrame)
	if len(timeFrameB) == 0 {
		timeFrameB = types.M1
	}

	return &SignalTile{
		IndicatorA: IndicatorSettings{
			Indicator:   indicatorA,
			RealTime:    raw.IndicatorA.RealTime,
			CandlesBack: raw.IndicatorA.CandleBack,
			TimeFrame:   timeFrameA,
		},
		IndicatorB: IndicatorSettings{
			Indicator:   indicatorB,
			RealTime:    raw.IndicatorB.RealTime,
			CandlesBack: raw.IndicatorB.CandleBack,
			TimeFrame:   timeFrameB,
		},
		Operand:     raw.Operand,
		Persistence: raw.Persistence,
	}, nil
}

func rawIndicatorToIndicator(raw rawIndicator) (types.Indicator, error) {
	refType := INDICATOR_REGISTRY[raw.IndicatorKey]

	if refType == nil {
		return nil, fmt.Errorf("indicator %s not found in mapping", raw.IndicatorKey)
	}

	indicator := reflect.Indirect(reflect.New(refType))

	for key, value := range raw.Data {
		key := strings.Title(key)

		field := indicator.FieldByName(key)
		if !field.IsValid() {
			return nil, fmt.Errorf("field %s is not found on %s", key, raw.IndicatorKey)
		}

		fieldType := field.Type().String()
		valueType := reflect.ValueOf(value).Type().String()

		if fieldType == "int" && valueType == "float64" {
			value = int(value.(float64))
		}

		if fieldType == "types.Indicator" {
			var nestedRawIndicator rawIndicator

			marshalled, err := json.Marshal(value)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(marshalled, &nestedRawIndicator)
			if err != nil {
				return nil, err
			}

			nestedIndicator, err := rawIndicatorToIndicator(nestedRawIndicator)
			if err != nil {
				return nil, err
			}

			value = nestedIndicator
		}

		if fieldType == "helpers.CandlePosition" {
			value = helpers.CandlePosition(value.(string))
		}

		field.Set(reflect.ValueOf(value))
	}

	if val, ok := indicator.Interface().(types.Indicator); ok {
		return val, nil
	}
	t := reflect.New(indicator.Type())
	t.Elem().Set(indicator)
	if val, ok := t.Interface().(types.Indicator); ok {
		return val, nil
	}

	return nil, fmt.Errorf("type assertion failed for %s", raw.IndicatorKey)
}

package types

type QueuedBacktest struct {
	Id string
}

type BacktestEventData struct {
	Type PositionHandlerEventType `json:"type"`
	Data interface{}              `json:"data"`
}

type BacktestEvent struct {
	Status    string             `json:"status"`
	EventData *BacktestEventData `json:"eventData,omitempty"`
}

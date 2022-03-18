package types

type QueuedBacktest struct {
	Id string
}

type BacktestEventData struct {
	Type PositionHandlerEventType
	Data interface{}
}

type BacktestEvent struct {
	Status    string             `json:"status"`
	EventData *BacktestEventData `json:"eventData,omitempty"`
}

export enum PositionState {
	OPENING = "OPENING",
	OPEN = "OPEN",
	CLOSING = "CLOSING",
	CLOSED = "CLOSED"
}

export enum OrderSide {
	BUY = "BUY",
	SELL = "SELL"
}

export interface OrderFill {
	time: Date
	rate: number
	quantity: number
	quoteFee: number
}

export interface Order {
	orderId: string
	time: Date
	side: OrderSide
	active: boolean
	size: number
	rate: number
	fills: OrderFill[]
}

export default interface Position {
  id?: string
  botId: string
  symbol: string
	state: PositionState
	openTime: Date
	closeTime?: Date
	orders: Order[]
}

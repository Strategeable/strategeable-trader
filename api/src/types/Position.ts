import { ObjectId } from "mongodb"

enum PositionState {
	OPENING = "OPENING",
	OPEN = "OPEN",
	CLOSING = "CLOSING",
	CLOSED = "CLOSED"
}

enum OrderSide {
	BUY = "BUY",
	SELL = "SELL"
}

interface OrderFill {
	time: Date
	rate: number
	quantity: number
	quoteFee: number
}

interface Order {
	orderId: string
	time: Date
	side: OrderSide
	active: boolean
	size: number
	rate: number
	fills: OrderFill[]
}

export default interface Position {
  id?: ObjectId
  botId: ObjectId
  symbol: string
	state: PositionState
	openTime: Date
	closeTime?: Date
	orders: Order[]
}

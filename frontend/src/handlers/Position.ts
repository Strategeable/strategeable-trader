import Position, { Order, PositionState } from '@/types/Position'

export default class PositionHandler {
  public id?: string
  public botId: string
  public symbol: string
  public state: PositionState
  public openTime: Date
  public closeTime?: Date
  public orders: Order[]

  constructor (position: Position) {
    this.id = position.id
    this.botId = position.botId
    this.symbol = position.symbol
    this.state = position.state
    this.openTime = position.openTime
    this.closeTime = position.closeTime
    this.orders = position.orders
  }

  getEntryQuoteSize (): number {
    let size = 0
    for (const order of this.orders) {
      const fees = order.fills.reduce((acc, curr) => curr.quoteFee + acc, 0)
      size += order.rate * order.size - fees
    }

    return size
  }

  getMaxDrawdown (): number {
    return -1.21
  }

  getResultIncludingFees (): number {
    return 0.1
  }
}

import { ObjectId } from "mongodb"

export type Exchange = 'binance' | 'kucoin'

export interface ExchangeConnection {
  _id?: string
  userId: ObjectId
  exchange: Exchange
  name: string
  createdOn: Date
  apiKey: string
  apiSecret: string
}

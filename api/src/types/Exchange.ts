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

export interface ExchangeBalance {
  asset: string
  amount: number
}

export interface Rate {
  exchange: string
  asset: string
  priceInBtc: number
  priceInEth: number
  priceInUsdt: number
}

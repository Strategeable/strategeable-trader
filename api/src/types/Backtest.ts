import { ObjectId } from "mongodb";

interface Trade {
  symbol: string
  amountIn: number
  amountOut: number
  entryPrice: number
  exitPrice: number
  entryDate: Date
  exitDate: Date
  fees: number
  buyPathName: string
  sellPathName: string
}

export default interface Backtest {
  _id?: ObjectId
  startedOn: Date
  strategy: any
  fromDate: Date
  toDate: Date
  startBalance: number
  endBalance?: number
  finished: boolean
  positions: any[]
}

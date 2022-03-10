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
  id?: ObjectId
  strategyId: ObjectId
  strategy: string
  startDate: Date
  endDate: Date
  startBalance: number
  endBalance: number
  exchange: string
  trades: Trade[]
}

import { Strategy } from "./Strategy"

export interface BacktestRequestParameters {
  strategyId: string
  fromDate: Date
  toDate: Date
  startBalance: number
}

interface PositionValue {
  date: string
  rate: number
  baseSize: number
  quoteFees: number
}

export interface Position {
  openedAt: string
  closedAt: string
  symbol: string
  entryValue: PositionValue
  exitValue: PositionValue
}

export interface BacktestResult {
  id: string
  status: string
  strategy: Strategy
  startedOn: Date
  startBalance: number
  fromDate: string
  toDate: string
  finished: boolean
  endBalance: number
  positions: Position[]
}

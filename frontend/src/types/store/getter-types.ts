import { State } from "@/store"
import { BacktestResult } from "../Backtest"
import Bot from "../Bot"
import { ExchangeConnection } from "../Exchange"
import { Theme } from "../general"
import { Strategy } from "../Strategy"

export type Getters = {
  loggedIn(state: State): boolean
  strategies(state: State): Strategy[]
  bots(state: State): Bot[]
  backtests(state: State): Record<string, BacktestResult[]>
  theme(state: State): Theme
  exchangeConnections(state: State): ExchangeConnection[]
}

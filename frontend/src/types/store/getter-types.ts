import { State } from "@/store"
import { Socket } from "socket.io-client"
import { BacktestResult } from "../Backtest"
import Bot from "../Bot"
import {  ExchangeBalance, ExchangeConnection, Rate } from "../Exchange"
import { Theme } from "../general"
import { Strategy } from "../Strategy"

export type Getters = {
  loggedIn(state: State): boolean
  strategies(state: State): Strategy[]
  bots(state: State): Bot[]
  backtests(state: State): BacktestResult[],
  backtestsByStrategy(state: State): (strategyId: string) => BacktestResult[];
  theme(state: State): Theme
  exchangeConnections(state: State): ExchangeConnection[],
  balances(state: State): ExchangeBalance[],
  rates(state: State): Rate[],
  socket(state: State): Socket | undefined,
  denominateIn(state: State): 'BTC' | 'ETH' | 'USD',
  getAssetRounding(state: State): (asset: string) => number
  hasUser(state: State): boolean
}

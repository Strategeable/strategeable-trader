import { State } from '@/store'
import { BacktestResult } from '../Backtest'
import Bot from '../Bot'
import { Exchange, ExchangeBalance, ExchangeConnection } from '../Exchange'
import { Theme } from '../general'
import { Strategy } from '../Strategy'

export enum MutationTypes {
  SET_JWT = 'SET_JWT',
  SET_STRATEGIES = 'SET_STRATEGIES',
  SET_STRATEGY = 'SET_STRATEGY',
  ADD_BACKTEST_RESULT = 'ADD_BACKTEST_RESULT',
  SET_BACKTESTS = 'SET_BACKTESTS',
  SET_THEME = 'SET_THEME',
  ADD_EXCHANGE_CONNECTION = 'ADD_EXCHANGE_CONNECTION',
  SET_EXCHANGE_CONNECTIONS = 'SET_EXCHANGE_CONNECTIONS',
  DELETE_EXCHANGE_CONNECTION = 'DELETE_EXCHANGE_CONNECTION',
  ADD_BOT = 'ADD_BOT',
  SET_BOTS = 'SET_BOTS',
  SET_BALANCES = 'SET_BALANCES',
  SET_RATE = 'SET_RATE'
}

export type Mutations<S = State> = {
  [MutationTypes.SET_JWT](state: S, payload: string): void
  [MutationTypes.SET_STRATEGIES](state: S, payload: Strategy[]): void
  [MutationTypes.SET_STRATEGY](state: S, payload: Strategy): void
  [MutationTypes.ADD_BACKTEST_RESULT](state: S, payload: BacktestResult): void
  [MutationTypes.SET_BACKTESTS](state: S, payload: { strategyId: string, backtests: BacktestResult[] }): void
  [MutationTypes.SET_THEME](state: S, payload: Theme): void
  [MutationTypes.ADD_EXCHANGE_CONNECTION](state: S, payload: ExchangeConnection): void
  [MutationTypes.SET_EXCHANGE_CONNECTIONS](state: S, payload: ExchangeConnection[]): void
  [MutationTypes.DELETE_EXCHANGE_CONNECTION](state: S, payload: string): void
  [MutationTypes.ADD_BOT](state: S, payload: Bot): void
  [MutationTypes.SET_BOTS](state: S, payload: Bot[]): void
  [MutationTypes.SET_BALANCES](state: S, payload: ExchangeBalance[]): void
  [MutationTypes.SET_RATE](state: S, payload: { exchange: Exchange, symbol: string, rate: number }): void
}

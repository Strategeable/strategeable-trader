import { State } from '@/store'
import { Socket } from 'socket.io-client'
import { BacktestResult } from '../Backtest'
import Bot from '../Bot'
import { Exchange, ExchangeBalance, ExchangeConnection, Rate } from '../Exchange'
import { Theme } from '../general'
import Position from '../Position'
import { Strategy } from '../Strategy'

export enum MutationTypes {
  IO_BACKTEST_EVENT = 'IO_BACKTEST_EVENT',
  SET_SOCKET = 'SET_SOCKET',
  SET_JWT = 'SET_JWT',
  SET_STRATEGIES = 'SET_STRATEGIES',
  SET_STRATEGY = 'SET_STRATEGY',
  ADD_BACKTEST_RESULT = 'ADD_BACKTEST_RESULT',
  HAS_USER = 'HAS_USER',
  DELETE_BACKTEST_RESULT = 'DELETE_BACKTEST_RESULT',
  ADD_BACKTEST_RESULTS = 'ADD_BACKTEST_RESULTS',
  SET_THEME = 'SET_THEME',
  ADD_EXCHANGE_CONNECTION = 'ADD_EXCHANGE_CONNECTION',
  SET_EXCHANGE_CONNECTIONS = 'SET_EXCHANGE_CONNECTIONS',
  DELETE_EXCHANGE_CONNECTION = 'DELETE_EXCHANGE_CONNECTION',
  ADD_BOT = 'ADD_BOT',
  SET_BOTS = 'SET_BOTS',
  ADD_POSITIONS = 'ADD_POSITIONS',
  SET_BALANCES = 'SET_BALANCES',
  SET_RATE = 'SET_RATE',
  SET_RATES = 'SET_RATES'
}

export type Mutations<S = State> = {
  [MutationTypes.IO_BACKTEST_EVENT](state: S, payload: any): void
  [MutationTypes.SET_SOCKET](state: S, payload: Socket): void
  [MutationTypes.SET_JWT](state: S, payload: string): void
  [MutationTypes.SET_STRATEGIES](state: S, payload: Strategy[]): void
  [MutationTypes.SET_STRATEGY](state: S, payload: Strategy): void
  [MutationTypes.DELETE_BACKTEST_RESULT](state: S, payload: string): void
  [MutationTypes.ADD_BACKTEST_RESULT](state: S, payload: BacktestResult): void
  [MutationTypes.ADD_BACKTEST_RESULTS](state: S, payload: BacktestResult[]): void
  [MutationTypes.SET_THEME](state: S, payload: Theme): void
  [MutationTypes.ADD_EXCHANGE_CONNECTION](state: S, payload: ExchangeConnection): void
  [MutationTypes.SET_EXCHANGE_CONNECTIONS](state: S, payload: ExchangeConnection[]): void
  [MutationTypes.DELETE_EXCHANGE_CONNECTION](state: S, payload: string): void
  [MutationTypes.ADD_BOT](state: S, payload: Bot): void
  [MutationTypes.SET_BOTS](state: S, payload: Bot[]): void
  [MutationTypes.ADD_POSITIONS](state: S, payload: Position[]): void
  [MutationTypes.SET_BALANCES](state: S, payload: ExchangeBalance[]): void
  [MutationTypes.SET_RATE](state: S, payload: { exchange: Exchange, asset: string, quoteAsset: string, rate: number }): void
  [MutationTypes.SET_RATES](state: S, payload: Rate[]): void
  [MutationTypes.HAS_USER](state: S, payload: boolean): void
}

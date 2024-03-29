import { State } from '@/store'
import { ActionContext } from 'vuex'
import { BacktestRequestParameters, BacktestResult } from '../Backtest'
import Bot, { LaunchParameters } from '../Bot'
import { Exchange, ExchangeBalance, ExchangeConnection, Rate } from '../Exchange'
import { Theme } from '../general'
import Position from '../Position'
import { Strategy } from '../Strategy'
import { Mutations } from './mutation-types'

export enum ActionTypes {
  INIT = 'INIT',
  CHANGE_COLOR_THEME = 'CHANGE_COLOR_THEME',
  CHECK_AUTH_STATE = 'CHECK_AUTH_STATE',
  LOGIN = 'LOGIN',
  REGISTER_ACCOUNT = 'REGISTER_ACCOUNT',
  LOAD_STRATEGIES = 'LOAD_STRATEGIES',
  LOAD_BOTS = 'LOAD_BOTS',
  LOAD_STRATEGY = 'LOAD_STRATEGY',
  SAVE_STRATEGY = 'SAVE_STRATEGY',
  RUN_BACKTEST = 'RUN_BACKTEST',
  STOP_BACKTEST = 'STOP_BACKTEST',
  LOAD_BACKTESTS = 'LOAD_BACKTESTS',
  LOAD_EXCHANGE_CONNECTIONS = 'LOAD_EXCHANGE_CONNECTIONS',
  ADD_EXCHANGE_CONNECTION = 'ADD_EXCHANGE_CONNECTION',
  DELETE_EXCHANGE_CONNECTION = 'DELETE_EXCHANGE_CONNECTION',
  LAUNCH_BOT = 'LAUNCH_BOT',
  LOAD_POSITIONS = 'LOAD_POSITIONS',
  LOAD_BALANCES = 'LOAD_BALANCES',
  LOAD_RATES = 'LOAD_RATES'
}

type AugmentedActionContext = {
  commit<K extends keyof Mutations>(
    key: K,
    payload: Parameters<Mutations[K]>[1]
  ): ReturnType<Mutations[K]>
} & Omit<ActionContext<State, State>, 'commit'>

export interface Actions {
  [ActionTypes.INIT]({ commit }: AugmentedActionContext): void
  [ActionTypes.CHANGE_COLOR_THEME]({ commit, state }: AugmentedActionContext, theme?: Theme): void
  [ActionTypes.LOGIN](
    { commit, dispatch }: AugmentedActionContext,
    details: { password: string }
  ): Promise<boolean>
  [ActionTypes.REGISTER_ACCOUNT](
    { commit }: AugmentedActionContext,
    details: { password: string }
  ): Promise<string | undefined>
  [ActionTypes.LOAD_STRATEGIES]({ commit }: AugmentedActionContext): void
  [ActionTypes.CHECK_AUTH_STATE]({ commit }: AugmentedActionContext): Promise<boolean>
  [ActionTypes.LOAD_BOTS]({ commit }: AugmentedActionContext): void
  [ActionTypes.LOAD_STRATEGY]({ commit, state }: AugmentedActionContext, id: string): Promise<Strategy | undefined>
  [ActionTypes.SAVE_STRATEGY]({ commit, state }: AugmentedActionContext, strategy: Strategy): Promise<string | undefined>
  [ActionTypes.RUN_BACKTEST](
    { commit }: AugmentedActionContext,
    backtestParams: BacktestRequestParameters
  ): Promise<BacktestResult | undefined>
  [ActionTypes.STOP_BACKTEST](
    { commit }: AugmentedActionContext,
    backtestId: string
  ): Promise<BacktestResult | undefined>
  [ActionTypes.LOAD_BACKTESTS]({ commit }: AugmentedActionContext, id: string): Promise<BacktestResult[]>
  [ActionTypes.LOAD_EXCHANGE_CONNECTIONS]({ commit }: AugmentedActionContext): Promise<ExchangeConnection[]>
  [ActionTypes.ADD_EXCHANGE_CONNECTION](
    { commit }: AugmentedActionContext,
    conn: ExchangeConnection
  ): Promise<{ error?: string, data?: ExchangeConnection }>
  [ActionTypes.DELETE_EXCHANGE_CONNECTION](
    { commit }: AugmentedActionContext,
    id: string
  ): Promise<{ error?: string, success?: boolean }>
  [ActionTypes.LAUNCH_BOT](
    { commit }: AugmentedActionContext,
    params: LaunchParameters
  ): Promise<{ error?: string, data?: Bot }>
  [ActionTypes.LOAD_POSITIONS](
    { commit }: AugmentedActionContext,
    open: boolean
  ): Promise<{ error?: string, data?: Position[] }>
  [ActionTypes.LOAD_BALANCES]({ commit }: AugmentedActionContext): Promise<ExchangeBalance[]>
  [ActionTypes.LOAD_RATES](
    { commit }: AugmentedActionContext,
    { exchange, coins }: { exchange: Exchange, coins: string[] }
  ): Promise<Rate[]>
}

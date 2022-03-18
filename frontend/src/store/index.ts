import { ActionTree, CommitOptions, createStore, DispatchOptions, GetterTree, MutationTree, Store } from 'vuex'
import axios from '@/helpers/axios'
import { Strategy } from '@/types/Strategy'
import { BacktestResult } from '@/types/Backtest'
import Bot from '@/types/Bot'

import { Exchange, ExchangeBalance, ExchangeConnection } from '@/types/Exchange'
import { Mutations, MutationTypes } from '@/types/store/mutation-types'
import { Actions, ActionTypes } from '@/types/store/action-types'
import { Theme } from '@/types/general'
import { Getters } from '@/types/store/getter-types'

export interface State {
  token: string | undefined
  strategies: Strategy[]
  bots: Bot[]
  backtestsByStrategyId: Record<string, BacktestResult[]>
  theme: 'light' | 'dark'
  exchangeConnections: ExchangeConnection[]
  balances: Record<Exchange, ExchangeBalance[]>
  rates: Record<string, number>
}

const getters: GetterTree<State, State> & Getters = {
  loggedIn: state => !!state.token,
  strategies: state => state.strategies,
  bots: state => state.bots,
  backtests: state => state.backtestsByStrategyId,
  theme: state => state.theme,
  exchangeConnections: state => state.exchangeConnections,
  rates: state => state.rates,
  balances: state => state.balances
}

const mutations: MutationTree<State> & Mutations = {
  [MutationTypes.SET_JWT] (state, token) {
    state.token = token
  },
  [MutationTypes.SET_STRATEGIES] (state, strategies) {
    state.strategies = strategies
  },
  [MutationTypes.SET_STRATEGY] (state, strategy) {
    if (!strategy.id) {
      state.strategies.push(strategy)
    } else {
      const existingIndex = state.strategies.findIndex(s => s.id === strategy.id)
      if (existingIndex !== -1) {
        state.strategies[existingIndex] = strategy
      } else {
        state.strategies.push(strategy)
      }
    }
  },
  [MutationTypes.ADD_BACKTEST_RESULT] (state, result) {
    let backtests = state.backtestsByStrategyId[result.strategy.id || '']
    if (!backtests) backtests = []
    backtests.push(result)
  },
  [MutationTypes.SET_BACKTESTS] (state, { strategyId, backtests }) {
    state.backtestsByStrategyId[strategyId] = backtests
  },
  [MutationTypes.SET_THEME] (state, theme) {
    state.theme = theme
  },
  [MutationTypes.ADD_EXCHANGE_CONNECTION] (state, conn) {
    state.exchangeConnections.push(conn)
  },
  [MutationTypes.SET_EXCHANGE_CONNECTIONS] (state, connections) {
    state.exchangeConnections = connections
  },
  [MutationTypes.DELETE_EXCHANGE_CONNECTION] (state, id) {
    state.exchangeConnections = state.exchangeConnections.filter(e => e.id !== id)
  },
  [MutationTypes.ADD_BOT] (state, bot) {
    state.bots.push(bot)
  },
  [MutationTypes.SET_BOTS] (state, bots) {
    state.bots = bots
  },
  [MutationTypes.SET_BALANCES] (state, balances) {
    state.balances = balances
  },
  [MutationTypes.SET_RATE] (state, { exchange, symbol, rate }) {
    state.rates[`${exchange}-${symbol}`] = rate
  }
}

const actions: ActionTree<State, State> & Actions = {
  [ActionTypes.INIT] ({ dispatch }) {
    dispatch(ActionTypes.LOAD_STRATEGIES)
    dispatch(ActionTypes.LOAD_BOTS)
    dispatch(ActionTypes.LOAD_EXCHANGE_CONNECTIONS)
    dispatch(ActionTypes.LOAD_BALANCES)
  },
  [ActionTypes.CHANGE_COLOR_THEME] ({ commit, state }, theme) {
    // Toggle the color theme between dark & light
    let newTheme: Theme = state.theme === 'dark' ? 'light' : 'dark'
    // If theme is given as an argument, always set that theme
    if (theme) newTheme = theme

    localStorage.setItem('theme', newTheme)

    const html: any = document.querySelector('html')
    html.classList.remove(state.theme)
    html.classList.add(newTheme)

    commit(MutationTypes.SET_THEME, newTheme)
  },
  async [ActionTypes.LOGIN] ({ commit, dispatch }, { username, password }) {
    try {
      const response = await axios.post('/login', { username, password })
      if (!response.data || !response.data.token) return false

      const { token } = response.data
      localStorage.setItem('jwt', token)
      dispatch(ActionTypes.INIT)
      commit(MutationTypes.SET_JWT, token)
      return true
    } catch (err) {
      console.error(err)
      return false
    }
  },
  async [ActionTypes.REGISTER_ACCOUNT] ({ commit }, { username, password }) {
    try {
      const response = await axios.post('/register', { username, password })
      if (!response.data) return 'Something went wrong'
      if (response.data.error) return response.data.error

      const { token } = response.data
      localStorage.setItem('jwt', token)
      commit(MutationTypes.SET_JWT, token)

      return undefined
    } catch (err) {
      console.error(err)
      return 'Something went wrong'
    }
  },
  async [ActionTypes.LOAD_STRATEGIES] ({ commit }) {
    try {
      const response = await axios.get('/strategy')
      if (!response.data) return

      commit(MutationTypes.SET_STRATEGIES, response.data)
    } catch (err) {
      console.error(err)
    }
  },
  async [ActionTypes.LOAD_BOTS] ({ commit }) {
    try {
      const response = await axios.get('/bot')
      if (!response.data) return

      commit(MutationTypes.SET_BOTS, response.data)
    } catch (err) {
      console.error(err)
    }
  },
  async [ActionTypes.LOAD_STRATEGY] ({ commit, state }, id) {
    const localStrategy = state.strategies.find(s => s.id === id)
    if (localStrategy) return localStrategy

    try {
      const response = await axios.get(`/strategy/${id}`)
      if (!response.data) return

      commit(MutationTypes.SET_STRATEGY, response.data)
      return response.data
    } catch (err) {
      console.error(err)
      return undefined
    }
  },
  async [ActionTypes.SAVE_STRATEGY] ({ commit }, strategy) {
    try {
      const response = strategy.id ? await axios.put('/strategy', { strategy }) : await axios.post('/strategy', { strategy })
      if (!response.data) return

      commit(MutationTypes.SET_STRATEGY, response.data)
      return response.data.id
    } catch (err) {
      console.error(err)
      return undefined
    }
  },
  async [ActionTypes.RUN_BACKTEST] ({ commit }, backtestParams) {
    try {
      const response = await axios.post('/backtest', backtestParams)

      // Start polling for the backtest result
      // TODO: this preferably shouldn't keep polling the API, websockets could potentially help
      while (true) {
        await new Promise(resolve => setTimeout(resolve, 2000))
        try {
          const backtestResponse = await axios.get(`/backtest/${response.data.backtestId}`)
          if (backtestResponse.status === 200) {
            if (backtestResponse.data.finished) {
              commit(MutationTypes.ADD_BACKTEST_RESULT, backtestResponse.data)
              return backtestResponse.data
            }
          }
        } catch (err) {
          continue
        }
      }
    } catch (err) {
      console.error(err)
      return undefined
    }
  },
  async [ActionTypes.LOAD_BACKTESTS] ({ commit }, strategyId) {
    if (!strategyId) return
    try {
      const response = await axios.get(`/backtest/strategy/${strategyId}`)
      commit(MutationTypes.SET_BACKTESTS, { strategyId, backtests: response.data })
      return response.data
    } catch (err) {
      console.error(err)
      return []
    }
  },
  async [ActionTypes.LOAD_EXCHANGE_CONNECTIONS] ({ commit }) {
    try {
      const response = await axios.get('/settings/exchange-connection')
      commit(MutationTypes.SET_EXCHANGE_CONNECTIONS, response.data)
      return response.data
    } catch (err) {
      console.error(err)
      return []
    }
  },
  async [ActionTypes.ADD_EXCHANGE_CONNECTION] ({ commit }, exchangeConnection: ExchangeConnection) {
    try {
      const response = await axios.post('/settings/exchange-connection', {
        exchange: exchangeConnection.exchange,
        name: exchangeConnection.name,
        apiKey: exchangeConnection.apiKey,
        apiSecret: exchangeConnection.apiSecret
      })

      if (response.status !== 200) return { error: 'Something went wrong' }

      commit(MutationTypes.ADD_EXCHANGE_CONNECTION, response.data)
      return { data: response.data }
    } catch (err: any) {
      if (err.response.status === 409) return { error: 'Name already exists' }
      return { error: 'Something went wrong' }
    }
  },
  async [ActionTypes.DELETE_EXCHANGE_CONNECTION] ({ commit }, id) {
    try {
      const response = await axios.delete(`/settings/exchange-connection/${id}`)
      if (response.status !== 200) return { error: 'Something went wrong' }

      commit(MutationTypes.DELETE_EXCHANGE_CONNECTION, id)
      return { success: true }
    } catch (err: any) {
      if (err.response.status === 409) return { error: 'Name already exists' }
      return { error: 'Something went wrong' }
    }
  },
  async [ActionTypes.LAUNCH_BOT] ({ commit }, params) {
    try {
      const response = await axios.post('/bot', params)
      if (response.status !== 200) return { error: 'Something went wrong' }

      commit(MutationTypes.ADD_BOT, response.data)
      return { data: response.data }
    } catch (err: any) {
      return { error: 'Something went wrong' }
    }
  },
  async [ActionTypes.LOAD_BALANCES] ({ commit }) {
    try {
      const response = await axios.get('/settings/balances')
      commit(MutationTypes.SET_BALANCES, response.data)
    } catch (err) {
      console.error(err)
    }
  }
}

const store = createStore<State>({
  state: {
    token: undefined,
    strategies: [],
    bots: [],
    backtestsByStrategyId: {},
    theme: 'dark',
    exchangeConnections: [],
    balances: {
      binance: [],
      kucoin: []
    },
    rates: {}
  },
  getters,
  mutations,
  actions,
  modules: {
  }
})

export default store

export type StoreType = Omit<
  Store<State>,
  'getters' | 'commit' | 'dispatch'
> & {
  commit<K extends keyof Mutations, P extends Parameters<Mutations[K]>[1]>(
    key: K,
    payload: P,
    options?: CommitOptions
  ): ReturnType<Mutations[K]>
} & {
  dispatch<K extends keyof Actions>(
    key: K,
    payload: Parameters<Actions[K]>[1],
    options?: DispatchOptions
  ): ReturnType<Actions[K]>
} & {
  getters: {
    [K in keyof Getters]: ReturnType<Getters[K]>
  }
}

export function useStore () {
  return store as StoreType
}

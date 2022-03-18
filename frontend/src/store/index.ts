import { ActionTree, CommitOptions, createStore, DispatchOptions, GetterTree, MutationTree, Store } from 'vuex'
import axios from '@/helpers/axios'
import { Strategy } from '@/types/Strategy'
import { BacktestResult } from '@/types/Backtest'
import Bot from '@/types/Bot'

import { ExchangeBalance, ExchangeConnection, Rate } from '@/types/Exchange'
import { Mutations, MutationTypes } from '@/types/store/mutation-types'
import { Actions, ActionTypes } from '@/types/store/action-types'
import { Theme } from '@/types/general'
import { Getters } from '@/types/store/getter-types'
import { Socket } from 'socket.io-client'

export interface State {
  token: string | undefined
  strategies: Strategy[]
  bots: Bot[]
  backtests: BacktestResult[]
  theme: 'light' | 'dark'
  exchangeConnections: ExchangeConnection[]
  balances: ExchangeBalance[]
  rates: Rate[]
  socket: Socket | undefined
  denominateIn: 'BTC' | 'ETH' | 'USD',
  assetRounding: Record<string, number>
}

const getters: GetterTree<State, State> & Getters = {
  loggedIn: state => !!state.token,
  strategies: state => state.strategies,
  bots: state => state.bots,
  backtests: state => state.backtests,
  backtestsByStrategy: state => strategyId => state.backtests.filter(b => b.strategy.id === strategyId),
  theme: state => state.theme,
  exchangeConnections: state => state.exchangeConnections,
  rates: state => state.rates,
  balances: state => state.balances,
  socket: state => state.socket,
  denominateIn: state => state.denominateIn,
  getAssetRounding: state => asset => state.assetRounding[asset] || 4
}

const mutations: MutationTree<State> & Mutations = {
  [MutationTypes.IO_BACKTEST_EVENT] (state, event) {
    const backtestId = event.id
    const backtest = state.backtests.find(b => b.id === backtestId)
    if (!backtest) return

    event.events.forEach((event: any) => {
      if (event.status === 'FINISHED') {
        backtest.finished = true
      }

      backtest.status = event.status

      if (!event.eventData) return

      if (event.eventData.type === 'POSITION_CLOSED') {
        backtest.positions.push(event.eventData.data)
      } else if (event.eventData.type === 'TOTAL_BALANCE_CHANGED') {
        backtest.endBalance = event.eventData.data
      }
    })
  },
  [MutationTypes.SET_SOCKET] (state, socket) {
    state.socket = socket
  },
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
  [MutationTypes.DELETE_BACKTEST_RESULT] (state, id) {
    state.backtests = state.backtests.filter(b => b.id !== id)
  },
  [MutationTypes.ADD_BACKTEST_RESULT] (state, result) {
    const exists = state.backtests.some(b => b.id === result.id)
    if (exists) return
    state.backtests.push(result)
  },
  [MutationTypes.ADD_BACKTEST_RESULTS] (state, backtests) {
    for (const backtest of backtests) {
      const exists = state.backtests.some(b => b.id === backtest.id)
      if (exists) continue

      state.backtests.push(backtest)
    }
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
  [MutationTypes.SET_RATES] (state, rates) {
    state.rates = rates
  },
  [MutationTypes.SET_RATE] (state, { exchange, asset, quoteAsset, rate }) {
    const rateObj = state.rates.find(r => r.exchange === exchange && r.asset === asset)
    if (rateObj) {
      rateObj.quote[quoteAsset] = rate
    } else {
      const quoteObj: Record<string, number> = {}
      quoteObj[quoteAsset] = rate

      state.rates.push({
        asset,
        exchange,
        quote: quoteObj
      })
    }
  }
}

const actions: ActionTree<State, State> & Actions = {
  async [ActionTypes.INIT] ({ dispatch }) {
    dispatch(ActionTypes.LOAD_STRATEGIES)
    dispatch(ActionTypes.LOAD_BOTS)
    dispatch(ActionTypes.LOAD_EXCHANGE_CONNECTIONS)
    const balances: ExchangeBalance[] = await dispatch(ActionTypes.LOAD_BALANCES)
    const mapping: Record<string, string[]> = {}
    balances.forEach(b => {
      if (mapping[b.exchange]) {
        mapping[b.exchange] = [...mapping[b.exchange], b.asset]
      } else {
        mapping[b.exchange] = [b.asset]
      }
    })

    for (const [exchange, coins] of Object.entries(mapping)) {
      dispatch(ActionTypes.LOAD_RATES, { exchange, coins: [...coins, 'BTC', 'ETH'] })
    }
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
      commit(MutationTypes.ADD_BACKTEST_RESULT, response.data)
      return response.data
    } catch (err) {
      console.error(err)
      return undefined
    }
  },
  async [ActionTypes.STOP_BACKTEST] ({ commit }, id) {
    try {
      const response = await axios.post(`/backtest/${id}/stop`)
      commit(MutationTypes.DELETE_BACKTEST_RESULT, id)
      return response.data
    } catch (err) {
      return err
    }
  },
  async [ActionTypes.LOAD_BACKTESTS] ({ commit }, strategyId) {
    if (!strategyId) return
    try {
      const response = await axios.get(`/backtest/strategy/${strategyId}`)

      commit(MutationTypes.ADD_BACKTEST_RESULTS, response.data)
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
        apiSecret: exchangeConnection.apiSecret,
        passPhrase: exchangeConnection.passPhrase
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
      return response.data
    } catch (err) {
      console.error(err)
      return []
    }
  },
  async [ActionTypes.LOAD_RATES] ({ commit }, { exchange, coins }) {
    try {
      const response = await axios.get(`/rates?exchange=${exchange}&coins=${coins.join(',')}`)
      if (response.status !== 200) return []

      commit(MutationTypes.SET_RATES, response.data)
      return response.data
    } catch (err) {
      console.error(err)
      return []
    }
  }
}

const store = createStore<State>({
  state: {
    token: undefined,
    strategies: [],
    bots: [],
    backtests: [],
    theme: 'dark',
    exchangeConnections: [],
    balances: [],
    rates: [],
    socket: undefined,
    denominateIn: 'BTC',
    assetRounding: {
      BTC: 6,
      USD: 2,
      USDT: 2,
      ETH: 4
    }
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

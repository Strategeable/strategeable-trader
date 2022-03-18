import { createStore, StoreOptions } from 'vuex'
import axios from '@/helpers/axios'
import { Strategy } from '@/types/Strategy'
import { BacktestResult } from '@/types/Backtest'
import Bot from '@/types/Bot'
import { ExchangeConnection } from '@/types/Exchange'

export default createStore({
  state: {
    token: undefined,
    strategies: [] as Strategy[],
    bots: [] as Bot[],
    backtestsByStrategyId: {} as Record<string, BacktestResult[]>,
    theme: 'dark',
    exchangeConnections: [] as ExchangeConnection[],
    socket: undefined
  },
  getters: {
    loggedIn: state => !!state.token,
    strategies: state => state.strategies,
    bots: state => state.bots,
    backtests: state => state.backtestsByStrategyId,
    theme: state => state.theme,
    exchangeConnections: state => state.exchangeConnections,
    socket: state => state.socket
  },
  mutations: {
    io_backtestEvent (state, data) {
      console.log(data)
    },
    SET_SOCKET (state, socket) {
      state.socket = socket
    },
    SET_JWT (state, token) {
      state.token = token
    },
    SET_STRATEGIES (state, strategies) {
      state.strategies = strategies
    },
    SET_STRATEGY (state, strategy) {
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
    ADD_BACKTEST_RESULT (state, result) {
      let backtests = state.backtestsByStrategyId[result.strategy._id]
      if (!backtests) backtests = []
      backtests.push(result)
    },
    SET_BACKTESTS (state, { strategyId, backtests }) {
      state.backtestsByStrategyId[strategyId] = backtests
    },
    SET_THEME (state, theme) {
      state.theme = theme
    },
    ADD_EXCHANGE_CONNECTION (state, conn) {
      state.exchangeConnections.push(conn)
    },
    SET_EXCHANGE_CONNECTIONS (state, connections) {
      state.exchangeConnections = connections
    },
    DELETE_EXCHANGE_CONNECTION (state, id) {
      state.exchangeConnections = state.exchangeConnections.filter(e => e.id !== id)
    },
    ADD_BOT (state, bot) {
      state.bots.push(bot)
    },
    SET_BOTS (state, bots) {
      state.bots = bots
    }
  },
  actions: {
    init ({ dispatch }) {
      dispatch('loadStrategies')
      dispatch('loadBots')
      dispatch('loadExchangeConnections')
    },
    changeColorTheme ({ commit, state }, theme) {
      // Toggle the color theme between dark & light
      let newTheme = state.theme === 'dark' ? 'light' : 'dark'
      // If theme is given as an argument, always set that theme
      if (theme) newTheme = theme

      localStorage.setItem('theme', newTheme)

      const html: any = document.querySelector('html')
      html.classList.remove(state.theme)
      html.classList.add(newTheme)

      commit('SET_THEME', newTheme)
    },
    async login ({ commit, dispatch, getters }, { username, password }) {
      try {
        const response = await axios.post('/auth/login', { username, password })
        if (!response.data || !response.data.token) return

        const { token } = response.data
        localStorage.setItem('jwt', token)

        getters.socket.emit('authorization', token)

        dispatch('init')
        commit('SET_JWT', token)
      } catch (err) {
        console.error(err)
      }
    },
    async registerAccount ({ commit }, { username, password }): Promise<string | undefined> {
      try {
        const response = await axios.post('/auth/register', { username, password })
        if (!response.data) return 'Something went wrong'
        if (response.data.error) return response.data.error

        const { token } = response.data
        localStorage.setItem('jwt', token)
        commit('SET_JWT', token)

        return undefined
      } catch (err) {
        console.error(err)
        return 'Something went wrong'
      }
    },
    async loadStrategies ({ commit }) {
      try {
        const response = await axios.get('/strategy')
        if (!response.data) return

        commit('SET_STRATEGIES', response.data)
      } catch (err) {
        console.error(err)
      }
    },
    async loadBots ({ commit }) {
      try {
        const response = await axios.get('/bot')
        if (!response.data) return

        commit('SET_BOTS', response.data)
      } catch (err) {
        console.error(err)
      }
    },
    async loadStrategy ({ commit, state }, id) {
      const localStrategy = state.strategies.find(s => s.id === id)
      if (localStrategy) return localStrategy

      try {
        const response = await axios.get(`/strategy/${id}`)
        if (!response.data) return

        commit('SET_STRATEGY', response.data)
        return response.data
      } catch (err) {
        console.error(err)
        return undefined
      }
    },
    async saveStrategy ({ commit }, strategy) {
      try {
        const response = strategy.id ? await axios.put('/strategy', { strategy }) : await axios.post('/strategy', { strategy })
        if (!response.data) return

        commit('SET_STRATEGY', response.data)
        return response.data.id
      } catch (err) {
        console.error(err)
        return undefined
      }
    },
    async runBacktest ({ commit }, backtestParams) {
      try {
        const response = await axios.post('/backtest', backtestParams)
        return response.data

        // Start polling for the backtest result
        // TODO: this preferably shouldn't keep polling the API, websockets could potentially help
        // while (true) {
        //   await new Promise(resolve => setTimeout(resolve, 2000))
        //   try {
        //     const backtestResponse = await axios.get(`/backtest/${response.data.backtestId}`)
        //     if (backtestResponse.status === 200) {
        //       if (backtestResponse.data.finished) {
        //         commit('ADD_BACKTEST_RESULT', backtestResponse.data)
        //         return backtestResponse.data
        //       }
        //     }
        //   } catch (err) {
        //     continue
        //   }
        // }
      } catch (err) {
        console.error(err)
        return undefined
      }
    },
    async loadBacktests ({ commit }, strategyId) {
      if (!strategyId) return
      try {
        const response = await axios.get(`/backtest/strategy/${strategyId}`)
        commit('SET_BACKTESTS', { strategyId, backtests: response.data })
      } catch (err) {
        console.error(err)
      }
    },
    async loadExchangeConnections ({ commit }) {
      try {
        const response = await axios.get('/settings/exchange-connection')
        commit('SET_EXCHANGE_CONNECTIONS', response.data)
      } catch (err) {
        console.error(err)
      }
    },
    async addExchangeConnection ({ commit }, exchangeConnection: ExchangeConnection) {
      try {
        const response = await axios.post('/settings/exchange-connection', {
          exchange: exchangeConnection.exchange,
          name: exchangeConnection.name,
          apiKey: exchangeConnection.apiKey,
          apiSecret: exchangeConnection.apiSecret
        })

        if (response.status !== 200) return { error: 'Something went wrong' }

        commit('ADD_EXCHANGE_CONNECTION', response.data)
        return { data: response.data }
      } catch (err: any) {
        if (err.response.status === 409) return { error: 'Name already exists' }
        return { error: 'Something went wrong' }
      }
    },
    async deleteExchangeConnection ({ commit }, id) {
      try {
        const response = await axios.delete(`/settings/exchange-connection/${id}`)
        if (response.status !== 200) return { error: 'Something went wrong' }

        commit('DELETE_EXCHANGE_CONNECTION', id)
        return { data: response.data }
      } catch (err: any) {
        if (err.response.status === 409) return { error: 'Name already exists' }
        return { error: 'Something went wrong' }
      }
    },
    async launchBot ({ commit }, params) {
      try {
        const response = await axios.post('/bot', params)
        if (response.status !== 200) return { error: 'Something went wrong' }

        commit('ADD_BOT', response.data)
        return { data: response.data }
      } catch (err: any) {
        return { error: 'Something went wrong' }
      }
    }
  },
  modules: {
  }
})

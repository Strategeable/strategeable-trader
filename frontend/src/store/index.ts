import { createStore } from 'vuex'
import axios from '@/helpers/axios'
import { Strategy } from '@/types/Strategy'
import { BacktestResult } from '@/types/Backtest'

export default createStore({
  state: {
    token: undefined,
    strategies: [] as Strategy[],
    backtestsByStrategyId: {} as Record<string, BacktestResult[]>,
    theme: 'dark'
  },
  getters: {
    loggedIn: state => !!state.token,
    strategies: state => state.strategies,
    backtests: state => state.backtestsByStrategyId,
    theme: state => state.theme
  },
  mutations: {
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
    }
  },
  actions: {
    changeColorTheme ({ commit, state }, theme) {
      let newTheme = state.theme === 'dark' ? 'light' : 'dark'
      if (theme) newTheme = theme

      localStorage.setItem('theme', newTheme)

      const html: any = document.querySelector('html')
      html.classList.remove(state.theme)
      html.classList.add(newTheme)

      commit('SET_THEME', newTheme)
    },
    async login ({ commit, dispatch }, { username, password }) {
      try {
        const response = await axios.post('/login', { username, password })
        if (!response.data || !response.data.token) return

        const { token } = response.data
        localStorage.setItem('jwt', token)
        dispatch('init')
        commit('SET_JWT', token)
      } catch (err) {
        console.error(err)
      }
    },
    async registerAccount ({ commit }, { username, password }): Promise<string | undefined> {
      try {
        const response = await axios.post('/register', { username, password })
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
    init ({ dispatch }) {
      dispatch('loadStrategies')
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
        while (true) {
          await new Promise(resolve => setTimeout(resolve, 2000))
          try {
            const backtestResponse = await axios.get(`/backtest/${response.data.backtestId}`)
            if (backtestResponse.status === 200) {
              if (backtestResponse.data.finished) {
                commit('ADD_BACKTEST_RESULT', backtestResponse.data)
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
    async loadBacktests ({ commit }, strategyId) {
      if (!strategyId) return
      try {
        const response = await axios.get(`/backtest/strategy/${strategyId}`)
        commit('SET_BACKTESTS', { strategyId, backtests: response.data })
      } catch (err) {
        console.error(err)
      }
    }
  },
  modules: {
  }
})

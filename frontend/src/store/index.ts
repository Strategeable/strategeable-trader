import { createStore } from 'vuex'
import axios from '@/helpers/axios'
import { Strategy } from '@/types/Strategy'

export default createStore({
  state: {
    token: undefined,
    strategies: [] as Strategy[]
  },
  getters: {
    loggedIn: state => !!state.token,
    strategies: state => state.strategies
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
    }
  },
  actions: {
    async login ({ commit }, { username, password }) {
      try {
        const response = await axios.post('/login', { username, password })
        if (!response.data || !response.data.token) return

        const { token } = response.data
        localStorage.setItem('jwt', token)
        commit('SET_JWT', token)
      } catch (err) {
        console.error(err)
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

        commit('SET_STRATEGY', strategy)
        return response.data.id
      } catch (err) {
        console.error(err)
        return undefined
      }
    }
  },
  modules: {
  }
})

import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Bots from '../views/Bots.vue'
import Bot from '../views/Bot.vue'
import Strategies from '../views/Strategies.vue'
import Strategy from '../views/Strategy.vue'
import Settings from '../views/Settings.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Bots
  },
  {
    path: '/bots',
    name: 'Bots',
    component: Bots
  },
  {
    path: '/bots/:id',
    name: 'Bot',
    component: Bot
  },
  {
    path: '/strategies',
    name: 'Strategies',
    component: Strategies
  },
  {
    path: '/strategies/:id',
    name: 'Strategy',
    component: Strategy
  },
  {
    path: '/settings',
    name: 'Settings',
    component: Settings
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Dashboard from '../views/Dashboard.vue'
import Bots from '../views/Bots.vue'
import Strategies from '../views/Strategies.vue'
import Strategy from '../views/Strategy.vue'
import Settings from '../views/Settings.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Dashboard',
    component: Dashboard
  },
  {
    path: '/bots',
    name: 'Bots',
    component: Bots
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

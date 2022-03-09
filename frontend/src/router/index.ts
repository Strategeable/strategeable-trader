import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router'
import Bots from '../views/Bots.vue'
import Strategies from '../views/Strategies.vue'
import Strategy from '../views/Strategy.vue'

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
    path: '/strategies',
    name: 'Strategies',
    component: Strategies
  },
  {
    path: '/strategies/:id',
    name: 'Strategy',
    component: Strategy
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router

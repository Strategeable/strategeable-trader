import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './assets/css/extra.scss'
import { IconDefinition, library } from '@fortawesome/fontawesome-svg-core'
import {
  faTimes, faAngleDown,
  faGreaterThan, faGreaterThanEqual,
  faLessThan, faLessThanEqual,
  faEquals, faNotEqual,
  faPencilAlt, faArrowUp,
  faArrowDown, faUndo,
  faRedo, faMoon,
  faSun, faCaretDown,
  faCaretUp, faTrash,
  faArrowRight, faArrowLeft,
  faSlidersH
} from '@fortawesome/free-solid-svg-icons'
import 'vue-select/dist/vue-select.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import vSelect from 'vue-select'
import VueSocketIo from 'vue-3-socket.io'
import socketIoClient from 'socket.io-client'
import { MutationTypes } from './types/store/mutation-types'

const icons: IconDefinition[] = [
  faTimes, faAngleDown, faGreaterThan, faGreaterThanEqual, faLessThan, faLessThanEqual,
  faEquals, faNotEqual, faPencilAlt, faArrowUp, faArrowDown, faUndo, faRedo, faMoon, faSun,
  faCaretDown, faCaretUp, faTrash, faArrowRight, faArrowLeft, faSlidersH
]

for (const icon of icons) {
  library.add(icon)
}

const socketio = new VueSocketIo({
  connection: socketIoClient('http://localhost:3000', {
    extraHeaders: {
      Authorization: `Bearer ${localStorage.getItem('jwt')}`
    },
    reconnection: true,
    reconnectionDelay: 2000,
    reconnectionAttempts: Infinity
  }),
  vuex: {
    store,
    actionPrefix: 'IO_',
    mutationPrefix: 'IO_'
  }
})

createApp(App)
  .use(store)
  .use(router)
  .use(socketio)
  .component('fa-icon', FontAwesomeIcon)
  .component('v-select', vSelect)
  .mount('#app')

store.commit(MutationTypes.SET_SOCKET, (socketio as any).io)

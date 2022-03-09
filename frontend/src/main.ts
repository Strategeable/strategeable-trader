import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import './assets/css/extra.scss'
import { library } from '@fortawesome/fontawesome-svg-core'
import {
  faTimes, faAngleDown,
  faGreaterThan, faGreaterThanEqual,
  faLessThan, faLessThanEqual,
  faEquals, faNotEqual,
  faPencilAlt, faArrowUp,
  faArrowDown, faUndo,
  faRedo
} from '@fortawesome/free-solid-svg-icons'
import 'vue-select/dist/vue-select.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import vSelect from 'vue-select'

library.add(faPencilAlt)
library.add(faTimes)
library.add(faAngleDown)
library.add(faGreaterThan)
library.add(faGreaterThanEqual)
library.add(faLessThan)
library.add(faLessThanEqual)
library.add(faEquals)
library.add(faNotEqual)
library.add(faArrowUp)
library.add(faArrowDown)
library.add(faUndo)
library.add(faRedo)

createApp(App)
  .use(store)
  .use(router)
  .component('fa-icon', FontAwesomeIcon)
  .component('v-select', vSelect)
  .mount('#app')

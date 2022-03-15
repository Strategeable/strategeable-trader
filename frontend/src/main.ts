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
  faCaretUp, faTrash
} from '@fortawesome/free-solid-svg-icons'
import 'vue-select/dist/vue-select.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import vSelect from 'vue-select'

const icons: IconDefinition[] = [
  faTimes, faAngleDown, faGreaterThan, faGreaterThanEqual, faLessThan, faLessThanEqual,
  faEquals, faNotEqual, faPencilAlt, faArrowUp, faArrowDown, faUndo, faRedo, faMoon, faSun,
  faCaretDown, faCaretUp, faTrash
]

for (const icon of icons) {
  library.add(icon)
}

createApp(App)
  .use(store)
  .use(router)
  .component('fa-icon', FontAwesomeIcon)
  .component('v-select', vSelect)
  .mount('#app')

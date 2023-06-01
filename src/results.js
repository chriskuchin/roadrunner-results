import { createApp } from 'vue'
import router from './router'
import App from './App.vue'

import { createPinia } from 'pinia'

import './css/results.scss'

/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

/* import specific icons */
import { faPlus, faStopwatch, faRuler, faRepeat, faPlay, faDownload, faFileCsv, faEllipsisV, faArrowLeftLong } from '@fortawesome/free-solid-svg-icons'

library.add(faPlus, faStopwatch, faRuler, faRepeat, faPlay, faDownload, faFileCsv, faEllipsisV, faArrowLeftLong)

var pinia = createPinia()
var app = createApp(App)
  .component('icon', FontAwesomeIcon)
  .use(router).use(pinia)

app.mount('#app')

require('./assets/')

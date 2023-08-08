import { createApp } from 'vue'
import router from './router'
import App from './App.vue'

import { createPinia } from 'pinia'

import './css/results.scss'

/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'


if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/service-worker.js').then(registration => {
      console.log('SW registered: ', registration);
    }).catch(registrationError => {
      console.log('SW registration failed: ', registrationError);
    });
  });
}

/* import specific icons */
import { faPlus, faStopwatch, faRuler, faRepeat, faPlay, faDownload, faFileCsv, faEllipsisV, faArrowLeftLong, faFlagCheckered } from '@fortawesome/free-solid-svg-icons'
library.add(faPlus, faStopwatch, faRuler, faRepeat, faPlay, faDownload, faFileCsv, faEllipsisV, faArrowLeftLong, faFlagCheckered)

var pinia = createPinia()
var app = createApp(App)
  .component('icon', FontAwesomeIcon)
  .use(router).use(pinia)

app.mount('#app')

require('./assets/')

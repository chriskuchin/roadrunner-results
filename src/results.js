import { createApp } from 'vue'
// import store from './store'
import router from './router'
import App from './App.vue'

import './css/results.scss'


var app = createApp(App).use(router)

app.mount('#app')
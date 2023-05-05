import { createApp } from 'vue'
import router from './router'
import App from './App.vue'

import { createPinia } from 'pinia'

import './css/results.scss'

var pinia = createPinia()
var app = createApp(App).use(router).use(pinia)

app.mount('#app')

require('./assets/')

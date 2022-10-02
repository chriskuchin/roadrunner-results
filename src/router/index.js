import { createRouter, createWebHashHistory } from 'vue-router'
import RaceTimer from '../views/RaceTimer.vue'
import RecordResults from '../views/RecordResults.vue'
import HomePage from '../views/HomePage.vue'

const routes = [
  { path: "/", component: HomePage },
  { path: "/timer", component: RaceTimer },
  { path: "/record", component: RecordResults },
]

const router = new createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
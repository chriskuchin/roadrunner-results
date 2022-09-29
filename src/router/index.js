import Vue from 'vue'
import VueRouter from 'vue-router'
import RaceTimer from '../views/RaceTimer.vue'
import RecordResults from '../views/RecordResults.vue'
import HomePage from '../views/HomePage.vue'
Vue.use(VueRouter)

const routes = [
  { path: "/", component: HomePage },
  { path: "/timer", component: RaceTimer },
  { path: "/record", component: RecordResults },
]

const router = new VueRouter({
  routes
})

export default router
import { createRouter, createWebHashHistory } from 'vue-router'

const HomePage = () => import('../views/HomePage.vue')
const RaceTimer = () => import('../views/RaceTimer.vue')
const RecordResults = () => import('../views/RecordResults.vue')
const LoginPage = () => import('../views/LoginPage.vue')
const SignupPage = () => import('../views/SignupPage.vue')

const routes = [
  { path: "/", component: HomePage },
  { path: "/timer", component: RaceTimer },
  { path: "/record", component: RecordResults },
  { path: "/login", component: LoginPage },
  { path: "/signup", component: SignupPage }
]

const router = new createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
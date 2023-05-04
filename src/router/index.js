import { createRouter, createWebHashHistory } from 'vue-router'

const HomePage = () => import('../views/HomePage.vue')
const RaceTimer = () => import('../views/RaceTimer.vue')
const RecordResults = () => import('../views/RecordResults.vue')
const LoginPage = () => import('../views/LoginPage.vue')
const SignupPage = () => import('../views/SignupPage.vue')
const RaceBasePage = () => import('../views/RaceBasePage.vue')
const RaceEventsPage = () => import('../views/RaceEventsPage.vue')
const RaceEventsResultsPage = () => import('../views/RaceEventsResultsPage.vue')

const routes = [
  { path: "/", component: HomePage },
  {
    path: "/races/:raceId", component: RaceBasePage, children: [
      { path: "events", component: RaceEventsPage },
      { path: "events/:eventId/record", component: RecordResults, name: "record" },
      { path: "events/:eventId/results", component: RaceEventsResultsPage },
      { path: "events/:eventId/timer", component: RaceTimer, name: "timer" }
    ]
  },
  { path: "/login", component: LoginPage },
  { path: "/signup", component: SignupPage }
]

const router = new createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
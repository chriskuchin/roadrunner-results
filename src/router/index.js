import { createRouter, createWebHashHistory } from 'vue-router'

const HomePage = () => import('../views/HomePage.vue')
const RaceTimer = () => import('../views/RaceTimer.vue')
const RecordResults = () => import('../views/RecordResults.vue')
const LoginPage = () => import('../views/LoginPage.vue')
const SignupPage = () => import('../views/SignupPage.vue')
const RegisterParticipant = () => import('../views/RegisterParticipant.vue')
const RacesPage = () => import('../views/RacesPage.vue')
const RacePage = () => import('../views/RacePage.vue')
const RaceBasePage = () => import('../views/RaceBasePage.vue')
const RaceEventsPage = () => import('../views/RaceEventsPage.vue')
const RaceParticipantsPage = () => import('../views/RaceParticipantsPage.vue')
const RaceResultsPage = () => import('../views/RaceResultsPage.vue')

const routes = [
  { path: "/", component: HomePage },
  { path: "/races", component: RacesPage, name: "races" },
  {
    path: "/races/:raceId", component: RaceBasePage, children: [
      { path: "", component: RacePage, name: "race" },
      { path: "timer", component: RaceTimer, name: "timer" },
      { path: "record", component: RecordResults, name: "record" },
      { path: "register", component: RegisterParticipant, name: "register" },
      { path: "events", component: RaceEventsPage },
      { path: "participants", component: RaceParticipantsPage },
      { path: "results", component: RaceResultsPage }
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
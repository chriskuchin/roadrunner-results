import { createRouter, createWebHistory } from "vue-router";

const HomePage = () => import("../views/HomePage.vue");
const RaceTimer = () => import("../views/RaceTimer.vue");
const RecordResults = () => import("../views/RecordResults.vue");
const HeatPage = () => import("../views/HeatPage.vue");
const LoginPage = () => import("../views/LoginPage.vue");
const SignupPage = () => import("../views/SignupPage.vue");
const RaceBasePage = () => import("../views/RaceBasePage.vue");
const RaceEventsPage = () => import("../views/RaceEventsPage.vue");
const RaceEventsResultsPage = () =>
	import("../views/RaceEventsResultsPage.vue");
const RaceParticipants = () => import("../views/RaceParticipants.vue");
const RaceDivisionResultsPage = () =>
	import("../views/RaceDivisionResultsPage.vue");
const RaceVolunteersPage = () => import("../views/RaceVolunteers.vue");
const RaceInfoPage = () => import("../views/RacePage.vue");
const AthletePage = () => import("../views/AthletePage.vue");
const DistanceResults = () => import("../views/DistanceResults.vue");

const routes = [
	{ path: "/", component: HomePage },
	{ path: "/athlete", component: AthletePage },
	{ path: "/heats", name: "heats", component: HeatPage },
	{
		path: "/races/:raceId",
		component: RaceBasePage,
		children: [
			{ path: "", name: "info", component: RaceInfoPage },
			{ path: "volunteers", name: "volunteers", component: RaceVolunteersPage },
			{
				path: "participants",
				name: "participants",
				component: RaceParticipants,
			},
			{
				path: "divisions",
				name: "divisions",
				component: RaceDivisionResultsPage,
			},
			{ path: "events", name: "events", component: RaceEventsPage },
			{
				path: "events/:eventId/record",
				component: RecordResults,
				name: "record",
			},
			{ path: "events/:eventId/results", component: RaceEventsResultsPage },
			{ path: "events/:eventId/timer", component: RaceTimer, name: "timer" },
			{ path: "events/:eventId/distance", component: DistanceResults, name: "distance" },
		],
	},
	{ path: "/login", component: LoginPage },
	{ path: "/signup", component: SignupPage },
];

const router = new createRouter({
	history: createWebHistory(),
	routes,
});

export default router;

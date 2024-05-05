import { createRouter, createWebHistory } from "vue-router";

const HomePage = () => import("../views/HomePage.vue");

const LoginPage = () => import("../views/Login.vue");
const SignupPage = () => import("../views/Signup.vue");

const RaceBasePage = () => import("../views/races/BasePage.vue");
const RaceVolunteersPage = () => import("../views/races/Volunteers.vue");

const Timer = () => import("../views/events/Timer.vue");

const RecordResults = () => import("../views/RecordResults.vue");
const RaceEventsPage = () => import("../views/RaceEventsPage.vue");
const RaceEventsResultsPage = () =>
	import("../views/RaceEventsResultsPage.vue");
const RaceParticipants = () => import("../views/RaceParticipants.vue");
const RaceDivisionResultsPage = () =>
	import("../views/RaceDivisionResultsPage.vue");
const RaceInfoPage = () => import("../views/RacePage.vue");
const AthletePage = () => import("../views/AthletePage.vue");
const DistanceResults = () => import("../views/DistanceResults.vue");
const DisplayBoard = () => import("../views/DisplayBoard.vue");
const LaneAssignment = () => import("../views/LaneAssignment.vue");

const routes = [
	{ path: "/", component: HomePage },
	{ path: "/athlete", component: AthletePage },
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
			{ path: 'events/:eventId/lanes', component: LaneAssignment, name: "lanes" },
			{
				path: "events/:eventId/record",
				component: RecordResults,
				name: "record",
			},
			{ path: "events/:eventId/results", component: RaceEventsResultsPage },
			{ path: "events/:eventId/timer", component: Timer, name: "timer" },
			{ path: "events/:eventId/distance", component: DistanceResults, name: "distance" },
			{ path: "events/:eventId/board", component: DisplayBoard, name: "display" }
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

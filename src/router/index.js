import { createRouter, createWebHistory } from "vue-router";

const HomePage = () => import("../views/HomePage.vue");

const LoginPage = () => import("../views/Login.vue");
const SignupPage = () => import("../views/Signup.vue");
const AthletePage = () => import("../views/AthletePage.vue");

const RaceBasePage = () => import("../views/races/BasePage.vue");
const RaceVolunteersPage = () => import("../views/races/Volunteers.vue");
const RaceEventsPage = () => import("../views/races/Events.vue");
const Participants = () => import("../views/races/Participants.vue");
const DivisionResults = () => import("../views/races/DivisionResults.vue");
const Overview = () => import("../views/races/Overview.vue");

const DisplayBoard = () => import("../views/events/DisplayBoard.vue");
const Timer = () => import("../views/events/Timer.vue");
const DistanceResults = () => import("../views/events/DistanceResults.vue");
const TimerResults = () => import("../views/events/TimerResults.vue");
const EventResults = () => import("../views/events/EventResults.vue");
const LaneAssignment = () => import("../views/events/LaneAssignment.vue");

const routes = [
	{ path: "/", component: HomePage },
	{ path: "/athlete", component: AthletePage },
	{
		path: "/races/:raceId",
		component: RaceBasePage,
		children: [
			{ path: "", name: "info", component: Overview },
			{ path: "volunteers", name: "volunteers", component: RaceVolunteersPage },
			{
				path: "participants",
				name: "participants",
				component: Participants,
			},
			{
				path: "divisions",
				name: "divisions",
				component: DivisionResults,
			},
			{ path: "events", name: "events", component: RaceEventsPage },
			{ path: 'events/:eventId/lanes', component: LaneAssignment, name: "lanes" },
			{
				path: "events/:eventId/record",
				component: TimerResults,
				name: "record",
			},
			{ path: "events/:eventId/results", component: EventResults },
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

import { defineStore } from "pinia";
import { setAuthHeader } from "../../api/auth";

export const useEventStore = defineStore("event", {
	state: () => ({
		timers: [],
		id: "",
		description: "",
		name: "",
		type: "",
		distance: 0,
	}),
	getters: {
		timerList: (state) => state.timers,
	},
	actions: {
		loadEvent(event) {
			this.id = event.eventId
			this.name = event.description
			this.type = event.type
			this.distance = event.distance
		},
		async loadEventByID(raceID, eventID) {
			if (this.id === eventID) {
				return
			}

			const url = `/api/v1/races/${raceID}/events/${eventID}`
			let res = await fetch(url, setAuthHeader({ method: "GET" }))
			if (!res.ok)
				return

			let event = await res.json()

			this.loadEvent(event)
		},
		async loadTimers(raceID, eventID) {
			const url = `/api/v1/races/${raceID}/events/${eventID}/timers`;
			this.timers = await (await fetch(url, { method: "GET" })).json();
		},
	},
});

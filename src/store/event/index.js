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
		eventPromise: null,
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

			if (this.eventPromise === null) {
				const url = `/api/v1/races/${raceID}/events/${eventID}`
				this.eventPromise = fetch(url, setAuthHeader({ method: "GET" }))
			}

			let res = await this.eventPromise
			this.eventPromise = null

			if (!res.ok) {
				return
			}

			try {
				let event = await res.json()
				this.loadEvent(event)
			} catch (_) { }

		},
		async loadTimers(raceID, eventID) {
			const url = `/api/v1/races/${raceID}/events/${eventID}/timers`;
			this.timers = await (await fetch(url, { method: "GET" })).json();
		},
	},
});

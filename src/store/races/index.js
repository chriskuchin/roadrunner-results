import { defineStore } from "pinia";

export const useRacesStore = defineStore("races", {
	state: () => ({
		races: [],
	}),
	getters: {
		getRaces: (state) => state.races,
	},
	actions: {
		async loadRaces() {
			this.races = await (
				await fetch("/api/v1/races", { method: "GET" })
			).json();
		},
	},
});

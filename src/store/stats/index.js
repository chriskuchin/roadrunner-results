import { defineStore } from "pinia";
import { setAuthHeader } from "../../api/auth";

export const useStatsStore = defineStore("stats", {
	state: () => ({
		id: "",
		name: "",
		ownerID: "",
		eventCount: 0,
		events: [],
		participantStats: {
			total: 0,
			finishers: 0,
			male: 0,
			female: 0,
			birthYearDistro: [],
		},
	}),
	getters: {
		getID: (state) => state.id,
		getName: (state) => state.name,
		yearLabels: (state) => [
			...new Set(state.participantStats.birthYearDistro.map((row) => row.year)),
		],
		totalParticipants: (state) => state.participantStats.total,
		totalFinishers: (state) => state.participantStats.finishers,
		eventTotal: (state) => state.eventCount,
		eventList: (state) => state.events,
		eventName: (state) => {
			return (eventID) => {
				const event = state.events.filter((val) => val.eventId === eventID)[0];

				if (event) return event.description;

				return "";
			};
		},
		maleValues() {
			const values = [];
			for (const year of this.yearLabels) {
				const val = this.participantStats.birthYearDistro
					.filter((val) => val.gender === "Male" && val.year === year)
					.map((val) => val.count);

				if (val.length > 0) values.push(val[0]);
				else values.push(0);
			}

			return values;
		},
		femaleValues() {
			const values = [];
			for (const year of this.yearLabels) {
				const val = this.participantStats.birthYearDistro
					.filter((val) => val.gender === "Female" && val.year === year)
					.map((val) => val.count);

				if (val.length > 0) values.push(val[0]);
				else values.push(0);
			}

			return values;
		},
		yearValues: (state) =>
			state.participantStats.birthYearDistro.map((row) => row.count),
	},
	actions: {
		async loadStats(id) {
			const race = await (
				await fetch(`/api/v1/races/${id}/stats`, await setAuthHeader({ method: "GET" }))
			).json();
			this.name = race.name;
			this.id = race.id;
			this.eventCount = race.event_count;
			this.participantStats.birthYearDistro =
				race.participant_stats.birth_year_distribution;
			this.participantStats.male = race.participant_stats.male;
			this.participantStats.female = race.participant_stats.female;
			this.participantStats.total = race.participant_stats.total;
			this.participantStats.finishers = race.participant_stats.finishers;
			this.eventCount = race.eventCount;
			this.events = race.events;
		},
	},
});

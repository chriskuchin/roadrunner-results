import { defineStore } from "pinia";
import { setAuthHeader } from "../../api/auth";
import { getEventResults } from "../../api/events";

export const useResultsStore = defineStore("results", {
	state: () => {
		return {
			teamOptions: new Set(),
			genderOptions: new Set(),
			yearOptions: new Set(),
			timerOptions: new Set(),
		};
	},
	getters: {
		teams: (state) => [...state.teamOptions].sort(),
		genders: (state) => [...state.genderOptions],
		years: (state) => [...state.yearOptions].sort(),
	},
	actions: {
		getResults: async function (
			raceID,
			eventID,
			name,
			gender,
			team,
			year,
			timers,
			order,
		) {
			const results = await getEventResults(
				raceID,
				eventID,
				name,
				gender,
				team,
				year,
				timers,
				order,
			);

			for (const elem of results) {
				this.yearOptions.add(elem.birth_year);
				this.teamOptions.add(elem.team);
				this.genderOptions.add(elem.gender);
				this.timerOptions.add(elem.timer_id);
			}
			return results;
		},
		recordRunnerResult: async (runner) => {
			const payload = {
				bib_number: runner.bib,
			};

			if (runner.timerId !== "" && runner.timerId !== "latest") {
				payload.timer_id = runner.timerId;
			}

			const url = `/api/v1/races/${runner.raceId}/events/${runner.eventId}/results`;

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "PUT",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(payload),
				}),
			);

			return res.ok;
		},
	},
});

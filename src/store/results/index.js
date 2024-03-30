import { defineStore } from "pinia";
import { getEventResults } from "../../api/events";
import { deleteResult, recordResult, updateResult } from "../../api/results";

export const useResultsStore = defineStore("results", {
	state: () => {
		return {
			results: [],
			raceId: "",
			eventId: "",
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
		loadResults: async function (
			raceId,
			eventId,
			name,
			gender,
			team,
			year,
			timers,
			order,
		) {
			this.raceId = raceId
			this.eventId = eventId

			this.results = await getEventResults(
				raceId,
				eventId,
				name,
				gender,
				team,
				year,
				timers,
				order,
			);

			for (const elem of this.results) {
				this.yearOptions.add(elem.birth_year);
				this.teamOptions.add(elem.team);
				this.genderOptions.add(elem.gender);
				this.timerOptions.add(elem.timer_id);
			}
			return this.results;
		},
		updateResultByRowId: async function(rowId, result, bib) {
			this.results[rowId]
			console.log(this.results[rowId].result_id, result, bib)

			await updateResult(this.raceId, this.eventId, this.results[rowId].result_id, result, bib)
		},
		getResultByRowId: function(rowId) {
			return this.results[rowId]
		},
		deleteResultByRowId: async function (rowId) {
			if (!this.raceId || !this.eventId)
				throw new Error("raceId or eventId not set", this.raceId, this.eventId)

			await deleteResult(this.raceId, this.eventId, this.results[rowId].result_id)
		},
		recordRunnerResult: async (runner) => {
			return await recordResult(runner.raceId, runner.eventId, runner.bib, runner.timerId)
		},
	},
});

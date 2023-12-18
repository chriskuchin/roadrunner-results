import { defineStore } from "pinia";
import { setAuthHeader } from "../../api/auth";

export const useDivisionsStore = defineStore("divisions", {
	state: () => ({
		divisions: [],
	}),
	getters: {},
	actions: {
		load: async function (raceID) {
			const url = `/api/v1/races/${raceID}/divisions`;
			const res = await fetch(url);

			if (res.ok) {
				this.divisions = await res.json();
			}
		},
		createDivision: async (raceID, description, genders, years) => {
			const url = `/api/v1/races/${raceID}/divisions`;

			const payload = {
				display: description,
				filters: [
					{
						key: "gender",
						values: genders,
					},
					{
						key: "birth_year",
						values: years,
					},
				],
			};

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify(payload),
				}),
			);

			console.log(payload, res.ok);
		},
	},
});

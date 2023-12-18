import { defineStore } from "pinia";

import { setAuthHeader } from "../../api/auth";

export const useParticipantsStore = defineStore("participants", {
	state: () => ({
		participants: [],
	}),
	getters: {
		currentParticipants: (state) => state.participants,
	},
	actions: {
		async updateParticipant(raceID, participantID, participant) {
			const url = `/api/v1/races/${raceID}/participants/${participantID}`;

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "PUT",
					body: JSON.stringify({
						first_name: participant.first_name,
						last_name: participant.last_name,
						bib_number: Number(participant.bib_number),
						gender: participant.gender,
						team: participant.team,
						birth_year: Number(participant.birth_year),
					}),
				}),
			);

			if (res.ok) {
				return true;
			}
			return false;
		},
		async uploadParticipantCSV(raceID, file) {
			const url = `/api/v1/races/${raceID}/participants/csv`;
			const formData = new FormData();
			formData.append("csv", file, file.name);

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "POST",
					body: formData,
				}),
			);

			if (res.ok) {
				return true;
			}
			return false;
		},
		async registerParticipant(
			raceID,
			firstName,
			lastName,
			bibNumber,
			gender,
			birthYear,
			team,
		) {
			// requires first_name, last_name, bib_number, gender, birth_year, team
			const url = `/api/v1/races/${raceID}/participants`;

			const res = await fetch(
				url,
				await setAuthHeader({
					method: "POST",
					headers: {
						"Content-Type": "application/json",
					},
					body: JSON.stringify({
						first_name: firstName,
						last_name: lastName,
						bib_number: Number(bibNumber),
						gender: gender,
						team,
						birth_year: Number(birthYear),
					}),
				}),
			);

			if (res.ok) {
				return true;
			}
			return false;
		},
		async loadParticipants(raceID, gender, team, year, name, limit, offset) {
			let url = `/api/v1/races/${raceID}/participants`;
			const params = {};

			if (gender && gender !== "") {
				params.gender = gender;
			}

			if (team && team !== "") {
				params.team = team;
			}

			if (year && year !== "") {
				params.year = year;
			}

			if (name && name !== "") {
				params.name = name;
			}

			params.limit = limit;
			params.offset = offset;

			url += `?${new URLSearchParams(params).toString()}`;
			this.participants = await (await fetch(url, { method: "GET" })).json();
		},
	},
});

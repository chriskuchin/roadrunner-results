import { defineStore } from "pinia";

import { setAuthHeader } from '../../api/auth'

export const useParticipantsStore = defineStore("participants", {
    state: () => ({
        participants: [],
    }),
    getters: {
        currentParticipants: (state) => state.participants
    },
    actions: {
        async registerParticipant(raceID, firstName, lastName, bibNumber, gender, birthYear, team) {
            // requires first_name, last_name, bib_number, gender, birth_year, team
            let url = `/api/v1/races/${raceID}/participants`

            let res = await fetch(url, await setAuthHeader({
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    first_name: firstName,
                    last_name: lastName,
                    bib_number: bibNumber,
                    gender: gender,
                    team, team,
                    birth_year: birthYear
                })
            }))

            if (res.ok) {
                return true
            }
            return false
        },
        async loadParticipants(raceID, gender, team, year, name, limit, offset) {
            let url = `/api/v1/races/${raceID}/participants`
            let params = {}

            if (gender && gender != "") {
                params['gender'] = gender
            }

            if (team && team != "") {
                params['team'] = team
            }

            if (year && year != "") {
                params['year'] = year
            }

            if (name && name != "") {
                params['name'] = name
            }

            params['limit'] = limit
            params['offset'] = offset

            url += "?" + new URLSearchParams(params).toString()
            this.participants = await (await fetch(url, { method: "GET" })).json()
        }
    }
})
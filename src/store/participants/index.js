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

            console.log(this.participants)
        }
    }
})
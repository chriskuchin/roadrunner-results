import { defineStore } from "pinia";

export const useEventStore = defineStore("event", {
    state: () => ({
        timers: [],
    }),
    getters: {
        timerList: state => state.timers
    },
    actions: {
        async loadTimers(raceID, eventID) {
            let url = "/api/v1/races/" + raceID + "/events/" + eventID + "/timers"
            this.timers = await (await fetch(url, { method: "GET" })).json()
        }
    }
})
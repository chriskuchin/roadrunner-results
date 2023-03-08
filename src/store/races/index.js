import { defineStore } from "pinia";

export const useRaceStore = defineStore("races", {
    state: () => ({
        id: "",
        name: "",
        ownerID: "",
        participantStats: {
            total: 0,
            male: 0,
            female: 0,
            birthYearDistro: []
        }
    }),
    getters: {
        getID: state => state.id,
        getName: state => state.name,
        yearLabels: state => state.participantStats.birthYearDistro.map(row => row.year),
        yearValues: state => state.participantStats.birthYearDistro.map(row => row.count),
    },
    actions: {
        async loadRace(id) {
            let race = await (await fetch("/api/v1/races/" + id, { method: "GET" })).json()
            this.name = race.name
            this.id = race.id
            this.participantStats.birthYearDistro = race.participant_stats.birth_year_distribution
            this.participantStats.male = race.participant_stats.male
            this.participantStats.female = race.participant_stats.female
            this.participantStats.total = race.participant_stats.total
            console.log("state:", race)
        }
    }
})
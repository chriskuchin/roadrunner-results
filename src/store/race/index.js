import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'

export const useRaceStore = defineStore("race", {
  state: () => ({
    id: "",
    name: "",
    ownerID: "",
    eventCount: 0,
    events: [],
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
    yearLabels: state => [...new Set(state.participantStats.birthYearDistro.map(row => row.year))],
    totalParticipants: state => state.participantStats.total,
    eventTotal: state => state.eventCount,
    eventList: state => state.events,
    eventName: state => {
      return function (eventID) {
        let event = state.events.filter(val => val.eventId == eventID)[0]

        if (event)
          return event.description

        return ""
      }
    },
    maleValues() {
      let values = []
      this.yearLabels.forEach(year => {
        let val = this.participantStats.birthYearDistro.filter(val => val.gender == "M" && val.year == year).map(val => val.count)

        if (val.length > 0)
          values.push(val[0])
        else
          values.push(0)
      })

      return values
    },
    femaleValues() {
      let values = []
      this.yearLabels.forEach(year => {
        let val = this.participantStats.birthYearDistro.filter(val => val.gender == "F" && val.year == year).map(val => val.count)

        if (val.length > 0)
          values.push(val[0])
        else
          values.push(0)
      })

      return values
    },
    yearValues: state => state.participantStats.birthYearDistro.map(row => row.count),
  },
  actions: {
    async listVolunteers(id) {
      let url = `/api/v1/races/${id}/volunteers`

      let res = await fetch(url, await setAuthHeader({}))

      if (res.ok) {
        return await res.json()
      }

      return []
    },
    async shareRace(id, emails) {
      let url = `/api/v1/races/${id}/volunteers`
      let payload = {
        emails: emails,
      }

      let res = await fetch(url, await setAuthHeader({
        method: "PUT",
        body: JSON.stringify(payload)
      }))

      if (res.ok) {
        // added volunteers
        // handle failed adds
      } else {
        // whole request failed
      }
    },
    async loadRace(id) {
      let race = await (await fetch("/api/v1/races/" + id, { method: "GET" })).json()
      this.name = race.name
      this.id = race.id
      this.eventCount = race.event_count
      this.participantStats.birthYearDistro = race.participant_stats.birth_year_distribution
      this.participantStats.male = race.participant_stats.male
      this.participantStats.female = race.participant_stats.female
      this.participantStats.total = race.participant_stats.total
      this.eventCount = race.eventCount
      this.events = race.events
    }
  }
})
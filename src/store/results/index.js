import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'

export const useResultsStore = defineStore("results", {
  state: () => {
    return {
      teamOptions: new Set(),
      genderOptions: new Set(),
      yearOptions: new Set(),
      timerOptions: new Set(),
    }
  },
  getters: {
    teams: (state) => [...state.teamOptions].sort(),
    genders: (state) => [...state.genderOptions],
    years: (state) => [...state.yearOptions].sort(),
  },
  actions: {
    getResults: async function (raceID, eventID, name, gender, team, year) {
      let url = `/api/v1/races/${raceID}/events/${eventID}/results`
      let filters = new URLSearchParams()

      if (name !== "") {
        filters.append("name", name)
      }

      if (gender.length > 0) {
        gender.forEach((gender) => filters.append("gender", gender))
      }

      if (team.length > 0) {
        team.forEach((team) => filters.append("team", team))
      }

      if (year.length > 0) {
        year.forEach((year) => filters.append("year", year))
      }

      let res = await fetch(url + "?" + filters.toString())

      if (!res.ok) {
        // this.handleError("Failed retrieving results")
        console.log("error")
        return []
      }
      else {
        var results = await res.json()
        results.forEach((element) => {
          this.yearOptions.add(element.birth_year)
          this.teamOptions.add(element.team)
          this.genderOptions.add(element.gender)
          this.timerOptions.add(element.timer_id)
        })
        return results
      }
    },
    recordRunnerResult: async function (runner) {
      let payload = {
        bib_number: runner.bib
      }

      if (runner.timerId != "") {
        payload.timer_id = runner.timerId
      }

      let url = "/api/v1/races/" + runner.raceId + "/events/" + runner.eventId + "/results"

      let res = await fetch(url, await setAuthHeader({
        method: "PUT",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(payload)
      }))

      // if (!res.ok)
      //   return []
    }
  }
})
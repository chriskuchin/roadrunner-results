import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'
import { getEventResults } from "../../api/events";

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
    getResults: async function (raceID, eventID, name, gender, team, year, timers) {
      var results = await getEventResults(raceID, eventID, name, gender, team, year, timers)
      results.forEach((element) => {
        this.yearOptions.add(element.birth_year)
        this.teamOptions.add(element.team)
        this.genderOptions.add(element.gender)
        this.timerOptions.add(element.timer_id)
      })
      return results
    },
    recordRunnerResult: async function (runner) {
      let payload = {
        bib_number: runner.bib
      }

      if (runner.timerId != "" && runner.timerId != "latest") {
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

      return res.ok
    }
  }
})
import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'

export const useResultsStore = defineStore("results", {
  state: () => ({
    results: {},
    finishers: 0,
    runners: 0
  }),
  getters: {
    getResults: function (state) {
      return state.results
    }
  },
  actions: {
    recordFinishTime: function (finishTime) {
    },
    recordRunnerResult: async function (runner) {
      let payload = {
        bib_number: runner.bib
      }

      if (runner.timerId != "") {
        payload.timer_id = runner.timerId
      }

      let url = "/api/v1/races/" + runner.raceId + "/events/" + runner.eventId + "/results"

      let res = await fetch(url, setAuthHeader({
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
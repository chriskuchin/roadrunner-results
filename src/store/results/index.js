import { defineStore } from "pinia";

export const useResultsStore = defineStore("results", {
  state: () => ({
    results: {}
  }),
  getters: {
    getResults: function (state) {
      return state.results
    }
  },
  actions: {
    recordFinishTime: function (finishTime) {
      if (!this.results[finishTime.place]) {
        this.results[finishTime.place] = {}
      }

      this.results[finishTime.place].finishTimestamp = finishTime.timestamp
      this.results[finishTime.place].finishMinutes = finishTime.minutes
      this.results[finishTime.place].finishSeconds = finishTime.seconds
      this.results[finishTime.place].finishMilliseconds = finishTime.milliseconds

    },
    recordRunnerResult: function (runner) {
      if (!this.results[runner.place]) {
        this.results[finishTime.place] = {}
      }

      this.results[runner.place].runnerTimestamp = runner.timestamp
      this.results[runner.place].runnerBib = runner.bib
    }
  }
})
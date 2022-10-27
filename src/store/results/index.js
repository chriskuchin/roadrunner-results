import { defineStore } from "pinia";

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
      this.finishers++

      this.results[this.finishers] = {}
      this.results[this.finishers].finishTimestamp = finishTime.timestamp
      this.results[this.finishers].finishMinutes = finishTime.minutes
      this.results[this.finishers].finishSeconds = finishTime.seconds
      this.results[this.finishers].finishMilliseconds = finishTime.milliseconds

    },
    recordRunnerResult: function (runner) {
      if (!this.results[++this.runners]) {
        this.results[++this.finishers] = {}
      }

      this.results[this.runners].runnerTimestamp = runner.timestamp
      this.results[this.runners].runnerBib = runner.bib
    }
  }
})
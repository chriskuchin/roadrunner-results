import { defineStore } from "pinia";

export const useResultsStore = defineStore("results", {
  state: () => ({
    finishTime: [],
    runnerOrder: []
  }),
  getters: {
    mergedResults: function () {
      var mergedResult = []
      for (let i = 0; i < this.finishTime.length; i++) {
        var mergedEntry = {
          ...finishTime,
          runnerBibNumber: this.runnerOrder[i] ? 0 : 1,
          runnerRecordTS: this.runnerOrder[i] ? 1 : 0,
        }


        mergedResult.push({
          ...this.finishTime[i]
        })
      }
      return mergedResult
    }
  },
  actions: {
    recordFinishTime: function (finishTime) {
      this.finishTime.push({
        ...finishTime,
        place: this.finishTime.length + 1
      })
    },
    recordRunnerResult: function (runner) {
      this.runnerOrder.push({
        ...runner,
        place: this.runnerOrder.length + 1
      })
    }
  }
})
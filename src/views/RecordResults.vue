<template>
  <div class="section">
    <div class="select is-small">
      <select @change="refreshData()" v-model="timerId">
        <option value="latest" selected>Latest ({{ results.length }})</option>
        <option v-for="(timer, index) in timers" :key="timer.id" :value="timer.id">
          Heat {{ index + 1 }} ({{ timer.count }})
        </option>
      </select>
    </div>
    <div class="tabs">
      <ul>
        <li @click="tabSelect('manual')" :class="{ 'is-active': isActiveTab('manual') }">
          <a>Manual</a>
        </li>
      </ul>
    </div>
    <div class="container">
      <result-input v-if="isActiveTab('manual')" :time="getFirstUnmatchedTime" :finisher="getFirstUnmatchedPlace"
        :total-results="getHeatTotalResults" :race-id="this.$route.params.raceId" :event-id="this.$route.params.eventId"
        :timer-id="this.timerId" @recorded-racer="refreshData()" />
    </div>
  </div>
  <div v-if="results.length > 0 && showResults">
    <table class="table" style="margin: 0px auto;">
      <thead>
        <th>Position</th>
        <th>Time</th>
        <th>Bib</th>
      </thead>
      <tbody>
        <tr v-for="(result, index) in results" :key="result.bib_number">
          <td>{{ index + 1 }}</td>
          <td>
            {{ formatMilliseconds(result.result_ms) }}
          </td>
          <td>{{ result.bib_number }}</td>
        </tr>
      </tbody>
    </table>
    <!-- <not :show="error.show" type="is-danger is-light" @close="dismissError">{{ error.msg }}</not> -->
  </div>
</template>

<script>
import { formatMilliseconds } from "../utilities";
import RacerInput from "../components/ResultsInput.vue";
import ResultsTable from "../components/ResultsTable.vue";
import Notification from '../components/Notification.vue';

export default {
  components: {
    "result-input": RacerInput,
    "results-table": ResultsTable,
    "not": Notification,
  },
  mounted: function () {
    this.refreshData()
  },
  unmounted: function () {
    clearTimeout(this.resultsRefresh)
  },
  data: function () {
    return {
      showResults: false,
      resultsRefresh: null,
      activeTab: "manual",
      timers: [],
      timerId: "latest",
      results: [],
      error: {
        show: false,
        msg: "",
      }
    };
  },
  methods: {
    refreshData: function () {
      clearTimeout(this.resultsRefresh)

      this.listTimers()
      this.getHeatResults()

      this.resultsRefresh = setTimeout(this.refreshData, 2500)
    },
    async getHeatResults() {
      let url = "/api/v1/races/" + this.$route.params.raceId + "/events/" + this.$route.params.eventId + "/results?timerId=" + this.timerId

      this.results = await (await fetch(url)).json()
    },
    async listTimers() {
      let res = await fetch("/api/v1/races/" + this.$route.params.raceId + "/events/" + this.$route.params.eventId + "/timers")

      this.timers = await res.json()
    },
    formatMilliseconds,
    isActiveTab: function (tab) {
      return this.activeTab == tab;
    },
    tabSelect: function (tab) {
      this.activeTab = tab;
    },
  },
  computed: {
    getFirstUnmatchedPlace: function () {
      for (let i = 0; i < this.results.length; i++) {
        let result = this.results[i]
        if (result.bib_number == "") {
          return i + 1
        }
      }
      return "Complete"
    },
    getFirstUnmatchedTime: function () {
      console.log(this.results)
      for (let i = 0; i < this.results.length; i++) {
        let result = this.results[i]
        if (result.bib_number == "") {
          return formatMilliseconds(result.result_ms)
        }
      }

      return "None"
    },
    getHeatTotalResults: function () {
      return this.results.length
    }
  },
};
</script>
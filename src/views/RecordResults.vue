<template>
  <div>
    <div class="section">
      <div class="select is-small">
        <select @change="this.getHeatResults" v-model="timerId">
          <option value="" selected>Latest</option>
          <option v-for="(timer, index) in timers" :key="timer.id" :value="timer.id">Heat {{ index + 1 }} ({{ timer.count
          }})
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
        <result-input v-if="isActiveTab('manual')" :race-id="this.$route.params.raceId"
          :event-id="this.$route.params.eventId" :timer-id="this.timerId" />
      </div>
    </div>
    {{ this.results }}
  </div>
</template>

<script>
import RacerInput from "../components/ResultsInput.vue";
import ResultsTable from "../components/ResultsTable.vue";

export default {
  components: {
    "result-input": RacerInput,
    "results-table": ResultsTable,
  },
  mounted: function () {
    this.listTimers()
  },
  data: function () {
    return {
      activeTab: "manual",
      timers: [],
      timerId: "",
      results: []
    };
  },
  methods: {
    async getHeatResults() {
      console.log("getHeatResults")
      let url = "/api/v1/races/" + this.$route.params.raceId + "/events/" + this.$route.params.eventId + "/results?timerId=" + this.timerId

      this.results = await (await fetch(url)).json()
      console.log(this.results)
    },
    async listTimers() {
      let res = await fetch("/api/v1/races/" + this.$route.params.raceId + "/events/" + this.$route.params.eventId + "/timers")

      this.timers = await res.json()
    },
    isActiveTab: function (tab) {
      return this.activeTab == tab;
    },
    tabSelect: function (tab) {
      this.activeTab = tab;
    },
  },
  computed: {
  },
};
</script>
<template>
  <div>
    <div class="field has-addons">
      <div class="select is-multiple is-small is-rounded">
        <select multiple size="1" v-model="heatFilter">
          <option v-for="heat in heats" :value="heat.id" :key="heat.id">{{ heat.description }}</option>
        </select>
      </div>
    </div>
    <div class="table-container">
      <table class="table">
        <thead>
          <th>Position</th>
          <th>Time</th>
          <th>Bib</th>
          <th>Debug</th>
        </thead>
        <tbody>
          <tr v-for="(result, place) in filteredResults" :key="place">
            <td>{{ place + 1 }}</td>
            <td>{{ formatMilliseconds(result.result_ms) }}</td>
            <td>{{ result.bib_number }}</td>
            <td>{{ result }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { useEventStore } from "../store/event";
import { mapStores } from "pinia";
import { formatMilliseconds } from "../utilities";

export default {
  data: function () {
    return {
      results: [],
      heats: [],
      heatFilter: [],
    };
  },
  mounted: function () {
    this.getResults()
    this.getHeats()
  },
  methods: {
    getResults: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/results`
      let res = await fetch(url)

      this.results = await res.json()
    },
    getHeats: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/timers`
      let res = await fetch(url)
      let timers = await res.json()

      timers.forEach((element, index) => {
        this.heats.push({
          description: `Heat ${index + 1}`,
          id: element.id,
          start: element.timer_start
        })
        // default select all heats
        this.heatFilter.push(element.id)
      })
    },
    formatMilliseconds,
  },
  computed: {
    ...mapStores(useEventStore),
    filteredResults: function () {
      let results = []
      this.results.forEach((element) => {
        if (this.heatFilter.includes(element.timer_id) || this.heatFilter.length == 0) {
          results.push(element)
        }
      })
      return results
    }
  },
};
</script>
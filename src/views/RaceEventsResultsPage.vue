<template>
  <div>
    <table class="table">
      <thead>
        <th>Position</th>
        <th>Time</th>
        <th>Bib</th>
        <th></th>
      </thead>
      <tbody>
        <tr v-for="(result, place) in results" :key="place">
          <td>{{ place }}</td>
          <td>
            {{ formatMilliseconds(result.result_ms) }}
          </td>
          <td>{{ result.bib_number }}</td>
          <td>{{ result }}</td>
        </tr>
      </tbody>
    </table>
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
    };
  },
  mounted: function () {
    this.getResults()
  },
  methods: {
    getResults: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/results`
      let res = await fetch(url)

      this.results = await res.json()
      console.log(this.results)
    },
    formatMilliseconds,
  },
  computed: {
    ...mapStores(useEventStore)
  },
};
</script>
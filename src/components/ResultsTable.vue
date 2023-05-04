<template>
  <div>
    <table class="table">
      <thead>
        <th>Position</th>
        <th>Time</th>
        <th>Bib</th>
        <th>Difference</th>
        <th></th>
      </thead>
      <tbody>
        <tr v-for="(result, place) in raceResults" :key="place">
          <td>{{ place }}</td>
          <td>
            {{ formatFinishTime(result) }}
          </td>
          <td>{{ result.runnerBib }}</td>
          <td>{{ recordingDifference(result) }}</td>
          <td>{{ result }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { mapState } from "pinia";
import { useResultsStore } from "../store/results";

export default {
  props: {},
  data: function () {
    return {};
  },
  methods: {
    formatFinishTime: function (result) {
    },
    recordingDifference: function (result) {
      if (result.runnerTimestamp && result.finishTimestamp)
        return result.runnerTimestamp - result.finishTimestamp;
      return 0;
    },
  },
  computed: {
    ...mapState(useResultsStore, { raceResults: "getResults" }),
  },
};
</script>
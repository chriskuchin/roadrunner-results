<template>
  <div class="container">
    {{ athlete.first_name }}
    {{ athlete.last_name }}
    {{ athlete.birth_year }}
    {{ athlete.gender }}

    <div class="table-container">
      <table class="table" style="min-width: 100%;">
        <thead>
          <th><abbr title="Event">evt</abbr></th>
          <th><abbr title="Distance (miles)">miles</abbr></th>
          <th><abbr title="Distance (kilometers)">kms</abbr></th>
          <th><abbr title="Pace (min/mile)">min/mile</abbr></th>
          <th><abbr title="Pace (min/km)">min/km</abbr></th>
          <th><abbr title="Time">time</abbr></th>
        </thead>
        <tbody>
          <tr v-for="result in results">
            <td>{{ result.event.description }}</td>
            <td>{{ calculateMiles(result.event.distance).toFixed(2) }}</td>
            <td>{{ calculateKilometers(result.event.distance).toFixed(2) }}</td>
            <td>{{ formatMilliseconds(calculatePerMilePace(result.result, result.event.distance)) }} min/mile</td>
            <td>{{ formatMilliseconds(calculatePerKPace(result.result, result.event.distance)) }} min/km</td>
            <td>{{ formatMilliseconds(result.result) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { setAuthHeader } from '../api/auth';
import { formatMilliseconds, calculatePerKPace, calculatePerMilePace, calculateKilometers, calculateMiles } from '../utilities';

export default {
  mounted: function () {
    if (this.$route.query.first_name)
      this.athlete.first_name = this.$route.query.first_name

    if (this.$route.query.last_name)
      this.athlete.last_name = this.$route.query.last_name

    if (this.$route.query.birth_year)
      this.athlete.birth_year = this.$route.query.birth_year

    if (this.$route.query.gender)
      this.athlete.gender = this.$route.query.gender

    this.getAthlete()
  },
  data: function () {
    return {
      athlete: {
        first_name: "",
        last_name: "",
        birth_year: 0,
        gender: ""
      },
      results: []
    }
  },
  computed: {
  },
  methods: {
    calculatePerKPace,
    calculateKilometers,
    calculatePerMilePace,
    calculateMiles,
    formatMilliseconds,
    getAthlete: async function () {
      let url = `/api/v1/athletes/results/search`
      let filters = new URLSearchParams()

      if (this.athlete.first_name !== "") {
        filters.append("first_name", this.athlete.first_name)
      }

      if (this.athlete.last_name !== "") {
        filters.append("last_name", this.athlete.last_name)
      }

      if (this.athlete.gender !== "") {
        filters.append("gender", this.athlete.gender)
      }

      if (this.athlete.birth_year !== "") {
        filters.append("birth_year", this.athlete.birth_year)
      }

      let res = await fetch(url + "?" + filters.toString(), await setAuthHeader({}))

      if (res.ok) {
        this.results = await res.json()
      }
    }
  }
};
</script>

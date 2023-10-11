<template>
  <div class="section">
    <div class="box" v-for="(division) in  divisionTables ">
      <h1 class="title">{{ division.display }}</h1>
      <div class="table-container">
        <table class="table" style="min-width: 100%;">
          <thead>
            <th>Position</th>
            <th>Time</th>
            <th>Bib</th>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Gender</th>
            <th>Birth Year</th>
            <th>Team</th>
          </thead>
          <tbody>
            <tr v-for="(result, place) in  division.results " :key="place">
              <td>{{ place + 1 }}</td>
              <td>{{ formatMilliseconds(result.result_ms) }}</td>
              <td><a :href="'https://alphapeak.io/events/2023_10_RegionXC/images/photos/' + result.bib_number + '.mp4'">{{
                result.bib_number }}</a></td>
              <td>{{ result.first_name }}</td>
              <td>{{ result.last_name }}</td>
              <td>{{ result.gender }}</td>
              <td>{{ result.birth_year }}</td>
              <td>{{ result.team }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from 'pinia';
import { useDivisionsStore } from '../store/divisions';
import { getEventResults } from '../api/events';
import { formatMilliseconds } from "../utilities";

export default {
  data: function () {
    return {
      results: {}
    }
  },
  mounted: async function () {
    await this.load(this.$route.params.raceId)
    this.divisions.forEach((division) => {
      var genders = []
      var birthYears = []
      division.filters.forEach((filter) => {
        if (filter.key == "gender") {
          genders = filter.values
        } else if (filter.key == "birth_year") {
          birthYears = filter.values
        }
      })
      getEventResults(this.$route.params.raceId, this.$route.params.eventId, "", genders, [], birthYears, []).then((results) => {
        if (results.length > 0)
          this.results[division.display] = results
      })
    })
  },
  methods: {
    ...mapActions(useDivisionsStore, ['load']),
    formatMilliseconds,
  },
  computed: {
    ...mapState(useDivisionsStore, ['divisions']),
    divisionTables: function () {
      let tables = []
      this.divisions.forEach((division) => {
        if (this.results[division.display] && this.results[division.display].length > 0) {
          tables.push({
            display: division.display,
            results: this.results[division.display],
          })
        }
      })
      return tables
    }
  }
}
</script>
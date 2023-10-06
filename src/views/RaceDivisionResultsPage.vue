<template>
  <div class="section">
    <div class="box" v-for="division in divisions">
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
            <tr v-for="(result, place) in results[division.display]" :key="place">
              <td>{{ place + 1 }}</td>
              <td>{{ formatMilliseconds(result.result_ms) }}</td>
              <td>{{ result.bib_number }}</td>
              <td>{{ result.first_name }}</td>
              <td>{{ result.last_name }}</td>
              <td>{{ result.gender }}</td>
              <td>{{ result.birth_year }}</td>
              <td>{{ result.team }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      {{ }}
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
        this.results[division.display] = results
      })
    })
  },
  methods: {
    ...mapActions(useDivisionsStore, ['load']),
    formatMilliseconds,
  },
  computed: {
    ...mapState(useDivisionsStore, ['divisions'])
  }
}
</script>
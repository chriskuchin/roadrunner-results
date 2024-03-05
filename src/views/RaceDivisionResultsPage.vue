<template>
  <div class="section">
    <div class="field">
      <div class="control">
        <label class="label is-large">Event</label>
        <div class="select is-large">
          <select v-model="eventId">
            <option v-for="event in events" :id="event.eventId" :value="event.eventId">{{ event.description }}</option>
          </select>
        </div>
      </div>

      <!-- <div class="control">
        <label class="label is-large">Division</label>
        <div class="select is-large">
          <select v-model="activeDivisions">
            <option v-for="division in divisions" :id="division" :value="division">{{ division.display }}</option>
          </select>
        </div>
      </div> -->
    </div>
    <div class="box" v-for="(division) in  divisionTables" :key="division">
      <h1 class="title">{{ division.display }}</h1>
      <div class="table-container">
        <table class="table" style="min-width: 100%;">
          <thead>
            <th>Position</th>
            <th>{{ getResultColumnHeader }}</th>
            <th>Bib</th>
            <th>First Name</th>
            <th>Last Name</th>
            <th>Gender</th>
            <th>Birth Year</th>
            <th>Team</th>
            <th>Other</th>
          </thead>
          <tbody>
            <tr v-for="(result, place) in  division.results " :key="place">
              <td>{{ place + 1 }}</td>
              <td>{{ formatResults(result.result_ms) }}</td>
              <td>
                <a :href="'https://alphapeak.io/events/2023_10_RegionXC/images/photos/' + result.bib_number + '.mp4'">
                  {{ result.bib_number }}
                </a>
              </td>
              <td>{{ result.first_name }}</td>
              <td>{{ result.last_name }}</td>
              <td>{{ result.gender }}</td>
              <td>{{ result.birth_year }}</td>
              <td>{{ result.team }}</td>
              <td>
                <router-link
                  :to="getResultsLink(result.first_name, result.last_name, result.birth_year, result.gender)">
                  results
                </router-link>
              </td>
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
import { useEventStore } from '../store/event';
import { getEventResults, getRaceEvents } from '../api/events';
import { formatMilliseconds, formatCentimeters } from "../utilities";

export default {
  data: function () {
    return {
      results: {},
      events: [],
      eventId: "",
      activeDivisions: {}
    }
  },
  mounted: async function () {
    await this.load(this.getRaceID())

    this.events = await getRaceEvents(this.getRaceID())
    if (this.$route.query.eventId && this.$route.query.eventId != "")
      this.eventId = this.$route.query.eventId
    else if (this.events.length > 0)
      this.eventId = this.events[0].eventId

    await this.loadEventByID(this.getRaceID(), this.eventId)
  },
  watch: {
    eventId(eventId) {
      console.log(eventId)
      this.results = {}
      this.$router.push({ path: this.$route.path, query: { eventId: eventId } })
      this.loadEventByID(this.getRaceID(), eventId).then(() => {
        let order = "asc"
        if (this.type === "distance")
          order = "desc"
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

          getEventResults(this.getRaceID(), eventId, "", genders, [], birthYears, [], order).then((results) => {
            if (results.length > 0)
              this.results[division.display] = results
          })
        })
      })
    },
  },
  methods: {
    ...mapActions(useDivisionsStore, ['load']),
    ...mapActions(useEventStore, ['loadEventByID']),
    formatResults(results) {
      if (this.type === "distance")
        return formatCentimeters(results, "ftin")
      else
        return formatMilliseconds(results)
    },
    getRaceID: function () {
      return this.$route.params.raceId
    },
    getResultsLink: function (first_name, last_name, birth_year, gender) {
      let url = `/athlete`
      let filters = new URLSearchParams()

      filters.append("first_name", first_name)
      filters.append("last_name", last_name)
      filters.append("gender", gender)
      filters.append("birth_year", birth_year)

      // let res = await fetch(, await setAuthHeader({}))
      return url + "?" + filters.toString()
    }
  },
  computed: {
    ...mapState(useEventStore, ['type']),
    ...mapState(useDivisionsStore, ['divisions']),
    getResultColumnHeader() {
      if (this.type === "distance")
        return "Distance"

      return "Time"
    },
    divisionTables: function () {
      let tables = []
      for (const division of this.divisions) {
        if (this.results[division.display] && this.results[division.display].length > 0) {
          tables.push({
            display: division.display,
            results: this.results[division.display],
          })
        }
      }
      return tables
    }
  }
}
</script>
<template>
  <div class="section">
    <div class="field is-horizontal">
      <div class="field-label">
        <label class="label">Filters:</label>
      </div>
      <div class="field-body">
        <div class="field has-addons" v-if="heats.length > 0">
          <div class="select is-multiple is-small">
            <select multiple v-model="heatFilter">
              <option v-for="heat in heats" :value="heat.id" :key="heat.id">{{ heat.description }}</option>
            </select>
          </div>
        </div>
        <div class="field has-addons">
          <div class="control">
            <input class="input" type="text" placeholder="Find by name" v-model="filters.name">
          </div>
          <div class="control">
            <a class="button is-info" @click="getResults">
              Search
            </a>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple is-small">
              <select v-model="filters.gender" @change="getResults" multiple>
                <option v-for="gender in options.genders" :key="gender" :value="gender">{{ gender }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple is-small">
              <select v-model="filters.team" @change="getResults" multiple>
                <option v-for="team in options.teams" :key="team" :value="team">{{ team }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple is-small">
              <select v-model="filters.year" @change="getResults" multiple>
                <option v-for="year in [...options.years].sort()" :key="year" :value="year">{{ year }}</option>
              </select>
            </div>
          </div>
        </div>
      </div>
    </div>
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
          <tr v-for="(result, place) in filteredResults" :key="place">
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
  </div>
</template>

<script>
import { useEventStore } from "../store/event";
import { mapStores, mapActions } from "pinia";
import { formatMilliseconds } from "../utilities";
import { useErrorBus } from "../store/error";

export default {
  data: function () {
    return {
      results: [],
      heats: [],
      heatFilter: [],
      photos: [],
      options: {
        years: new Set(),
        teams: new Set(),
        genders: new Set(),
      },
      filters: {
        name: "",
        gender: [],
        team: [],
        year: [],
      }
    };
  },
  mounted: function () {
    this.getResults()
    this.getHeats()
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    getResults: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/results`
      let filters = new URLSearchParams()

      if (this.filters.name !== "") {
        filters.append("name", this.filters.name)
      }

      if (this.filters.gender.length > 0) {
        this.filters.gender.forEach((gender) => filters.append("gender", gender))
      }

      if (this.filters.team.length > 0) {
        this.filters.team.forEach((team) => filters.append("team", team))
      }

      if (this.filters.year.length > 0) {
        this.filters.year.forEach((year) => filters.append("year", year))
      }

      let res = await fetch(url + "?" + filters.toString())

      if (!res.ok)
        this.handleError("Failed retrieving results")
      else {
        this.results = await res.json()
        var that = this
        this.results.forEach((element) => {
          that.options.years.add(element.birth_year)
          that.options.teams.add(element.team)
          that.options.genders.add(element.gender)
        })

      }
    },
    getImageKeys: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/results/photos`
      let res = await fetch(url)

      if (!res.ok) {
        this.handleError("Failed Retrieving Photo Finishes")
      } else {
        this.photos = await res.json()
      }
    },
    getImageSrc(place) {
      if (this.photos[place]) {
        return `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/results/photos/${this.photos[place]}`
      }

      return ""
    },
    getHeats: async function () {
      let url = `/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/timers`
      let res = await fetch(url)

      if (!res.ok)
        this.handleError("Failed retrieving heats")
      else {
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
      }
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
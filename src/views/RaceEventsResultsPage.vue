<template>
  <div class="section">
    <div class="field is-horizontal">
      <div class="field-body">
        <div class="field has-addons">
          <div class="control">
            <input class="input" type="text" placeholder="Find by name" v-model="filters.name">
          </div>
          <div class="control">
            <a class="button is-info" @click="calculateResults">
              Search
            </a>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select">
              <select @change="selectDivision">
                <option selected value="0">Select a Division</option>
                <option v-for="(division, key) in getDivisions()" :key="key" :value="division">{{ key }}
                </option>
              </select>
            </div>
          </div>
        </div>
        <div class="field is-narrow" v-if="heats.length > 0">
          <div class="select is-multiple">
            <select multiple v-model="heatFilter">
              <option v-for="heat in heats" :value="heat.id" :key="heat.id">{{ heat.description }}</option>
            </select>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple">
              <select v-model="filters.gender" multiple>
                <option v-for="gender in genders" :key="gender" :value="gender">{{ gender }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple">
              <select v-model="filters.team" multiple>
                <option v-for="team in teams" :key="team" :value="team">{{ team }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="field is-narrow">
          <div class="control">
            <div class="select is-multiple">
              <select v-model="filters.year" multiple>
                <option v-for="year in years" :key="year" :value="year">{{ year }}</option>
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
          <th>Division</th>
        </thead>
        <tbody>
          <tr v-for="(result, place) in calculatedResults" :key="place">
            <td>{{ place + 1 }}</td>
            <td>{{ formatMilliseconds(result.result_ms) }}</td>
            <td>{{ result.bib_number }}</td>
            <td>{{ result.first_name }}</td>
            <td>{{ result.last_name }}</td>
            <td>{{ result.gender }}</td>
            <td>{{ result.birth_year }}</td>
            <td>{{ result.team }}</td>
            <td>{{ getParticipantDivision(result.birth_year) }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import { useEventStore } from "../store/event";
import { useResultsStore } from "../store/results"
import { mapStores, mapActions, mapState } from "pinia";
import { formatMilliseconds } from "../utilities";
import { useErrorBus } from "../store/error";

export default {
  data: function () {
    return {
      calculatedResults: [],
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
    this.calculateResults()
    this.getHeats()
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    ...mapActions(useResultsStore, ['getResults']),
    getDivisions: function () {
      let year = new Date().getFullYear()
      let firstDivision = `${year - 7}+`
      let yougestDivisionFilter = []
      for (var i = year - 7; i < year; i++) {
        yougestDivisionFilter.push(i)
      }
      let divisions = {}
      divisions[firstDivision] = yougestDivisionFilter

      for (var i = 1; i <= 10; i = i + 2) {
        var high = year - (7 + i)
        var low = year - (7 + i + 1)
        var currentDivision = `${low}-${high}`

        divisions[currentDivision] = [low, high]
      }

      return divisions
    },
    getParticipantDivision: function (birthYear) {
      let divisions = this.getDivisions()
      let participantDivision = "Unknown"

      Object.keys(divisions).forEach((key) => {
        if (divisions[key].includes(birthYear)) {
          participantDivision = key
          return
        }
      })

      return participantDivision
    },
    selectDivision: function (evt) {
      let years = evt.target.value.split(",")
      this.filters.year = []
      years.forEach(year => this.filters.year.push(year))
      evt.target.value = "0"
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
    calculateResults: async function () {
      this.calculatedResults = await this.getResults(this.$route.params.raceId, this.$route.params.eventId, this.filters.name, this.filters.gender, this.filters.team, this.filters.year)
    },
  },
  watch: {
    filters: {
      deep: true,
      handler() {
        this.calculateResults()
      },
    },
  },
  computed: {
    ...mapState(useResultsStore, ['years', 'genders', 'teams']),
    ...mapStores(useEventStore),
  },
};
</script>
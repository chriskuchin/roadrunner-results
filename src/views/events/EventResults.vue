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
          <div class="select is-multiple">
            <select multiple v-model="filters.timers">
              <option v-for="heat in options.heats" :value="heat.id" :key="heat.id">{{ heat.description }}
              </option>
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
          <th>Result</th>
          <th>Bib</th>
          <th>First Name</th>
          <th>Last Name</th>
          <th>Gender</th>
          <th>Birth Year</th>
          <th>Team</th>
          <th>Division</th>
          <th v-if="isLoggedIn"></th>
        </thead>
        <tbody class="is-size-4">
          <tr v-for="(result, rowid) in results" :key="rowid" class="is-size-4">
            <td>{{ rowid + 1 }}</td>
            <td>{{ formatResults(result.result_ms) }}</td>
            <td>{{ result.bib_number }}</td>
            <td>{{ result.first_name }}</td>
            <td>{{ result.last_name }}</td>
            <td>{{ result.gender }}</td>
            <td>{{ result.birth_year }}</td>
            <td>{{ result.team }}</td>
            <td>{{ getParticipantDivision(result.birth_year) }}</td>
            <td v-if="isLoggedIn">
              <div class="field has-addons">
                <p class="control">
                  <button class="button is-small is-info" @click="editResult(rowid)">
                    <span class="icon is-small">
                      <icon icon="fa-solid fa-pencil"></icon>
                    </span>
                  </button>
                </p>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <edit-results ref="edit-results" @reload="calculateResults" />
    <!-- <fab>
      <icon icon="fa-solid fa-plus"></icon>
    </fab> -->
  </div>
</template>

<script>
import { useEventStore } from "../../store/event";
import { useResultsStore } from "../../store/results"
import { mapActions, mapState } from "pinia";
import { useErrorBus } from "../../store/error";
import { useUserStore } from "../../store/user";
import { formatMilliseconds, formatCentimeters } from "../../utilities";
import EditResultsModal from '../../components/modals/EditResults.vue'
import Fab from "../../components/Fab.vue";

export default {
  components: {
    'edit-results': EditResultsModal,
    'fab': Fab,
  },
  data: function () {
    return {
      photos: [],
      options: {
        heats: [],
      },
      filters: {
        name: "",
        gender: [],
        team: [],
        year: [],
        timers: [],
      }
    };
  },
  mounted: function () {
    const that = this
    this.loadEventByID(this.$route.params.raceId, this.$route.params.eventId).then(() => {
      that.calculateResults()
    })
    this.getHeats()
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    ...mapActions(useResultsStore, ['loadResults']),
    ...mapActions(useEventStore, ['loadEventByID']),
    editResult: function (rowId) {
      this.$refs['edit-results'].open(rowId)
    },
    formatResults: function (result) {
      if (this.type === "distance")
        return formatCentimeters(result, "ftin")
      else
        return formatMilliseconds(result)
    },
    getDivisions: function () {
      let year = new Date().getFullYear()
      let firstDivision = `${year - 6}+`
      let yougestDivisionFilter = []
      for (var i = year - 6; i < year; i++) {
        yougestDivisionFilter.push(i)
      }
      let divisions = {}
      divisions[firstDivision] = yougestDivisionFilter

      for (var i = 1; i <= 10; i = i + 2) {
        var high = year - (6 + i)
        var low = year - (6 + i + 1)
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
      let res = await fetch(`/api/v1/races/${this.$route.params.raceId}/events/${this.$route.params.eventId}/timers`)

      if (!res.ok)
        this.handleError("Failed retrieving heats")
      else {
        let timers = await res.json()

        if (timers == null) {
          timers = []
        }

        timers.forEach((element, index) => {
          this.options.heats.push({
            description: `Heat ${index + 1}`,
            id: element.id,
            start: element.timer_start
          })
          // default select all heats
          this.filters.timers.push(element.id)
        })
      }
    },
    formatCentimeters,
    formatMilliseconds,
    calculateResults: async function () {
      let order = "asc"
      if (this.type === "distance")
        order = "desc"

      this.loadResults(this.$route.params.raceId, this.$route.params.eventId, this.filters.name, this.filters.gender, this.filters.team, this.filters.year, this.filters.timers, order)
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
    ...mapState(useUserStore, ['isLoggedIn']),
    ...mapState(useResultsStore, ['years', 'genders', 'teams', 'results']),
    ...mapState(useEventStore, ['type'])
  },
};
</script>
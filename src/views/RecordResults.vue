<template>
  <div class="mx-auto">
    <div class="select is-small">
      <select @change="refreshData" v-model="timerId">
        <option v-for="(timer, index) in timers" :key="timer.id" :value="timer.id">
          Heat {{ index + 1 }} ({{ timer.count }})
        </option>
      </select>
    </div>
    <div class="tabs mt-4">
      <ul>
        <li @click="tabSelect('manual')" :class="{ 'is-active': isActiveTab('manual') }">
          <a>Manual</a>
        </li>
        <li v-if="hasLaneAssignments" @click="tabSelect('heat')" :class="{ 'is-active': isActiveTab('heat') }">
          <a>Heat</a>
        </li>
        <li @click="tabSelect('scan')" :class="{ 'is-active': isActiveTab('scan') }">
          <a>Scanner</a>
        </li>
      </ul>
    </div>
    <div class="container">
      <div class="mx-2 mt-4" v-if="isActiveTab('manual')">
        <div class="columns">
          <result-input class="column recorder" :time="getFirstUnmatchedTime" :finisher="getFirstUnmatchedPlace"
            :total-results="getHeatTotalResults" :race-id="this.$route.params.raceId"
            :event-id="this.$route.params.eventId" :timer-id="this.timerId" @bib="bibInput" />
          <div class="column results">
            <div class="table-container">
              <tbl class="mx-auto is-narrow" :headers="heatResultsHeader" :rows="heatResults" />
            </div>
          </div>
        </div>
      </div>
      <scan v-else-if="isActiveTab('scan')" @bib="bibInput" />
      <div class="mx-2 mt-4" v-else-if="isActiveTab('heat')">
        <div class="columns">
          <div class="column recorder">
            <div class="fixed-grid has-3-cols mx-auto">
              <div class="grid">
                <div class="cell" v-for="assignment in laneAssignments">
                  <button class="button is-large is-primary is-fullwidth" @click="recordLaneFinish(assignment)"
                    :disabled="assignment.bib === ''">
                    {{ assignment.lane }}
                  </button>
                </div>
              </div>
            </div>
            <div class="buttons is-right">
              <button class="button is-link" @click="saveHeat"
                :disabled="results.length != heatFinish.length">Save</button>
            </div>
          </div>
          <div class="column">
            <div class="table-container">
              <tbl class="mx-auto is-narrow" :headers="heatResultsHeader" :rows="heatResults" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { formatMilliseconds } from "../utilities";
import { recordResult, getHeatResults } from "../api/results";
import { listTimers } from '../api/timers';
import { getParticipantByBib } from "../api/participants";
import { mapState, mapActions } from "pinia";
import { useParticipantsStore } from '../store/participants';
import RacerInput from "../components/ResultsInput.vue";
import Scanner from '../components/Scanner.vue';
import Table from '../components/Table.vue';

export default {
  components: {
    "result-input": RacerInput,
    "scan": Scanner,
    "tbl": Table,
  },
  mounted: function () {
    this.loadParticipants(this.$route.params.raceId, "", "", "", "", 500, 0)
    this.refreshData()
  },
  unmounted: function () {
    clearTimeout(this.resultsRefresh)
  },
  data: function () {
    return {
      resultsRefresh: null,
      activeTab: "manual",
      timers: [],
      timerId: "latest",
      results: [],
      heatFinish: [],
      heatFinisher: new Map(),
      heatResultsHeader: [
        {
          abbr: "Pos",
          title: "Position",
        },
        {
          abbr: "Time",
          title: "Finished Time"
        },
        {
          abbr: "Ln",
          title: "Lane",
        },
        {
          abbr: "Bib",
          title: "Bib Number"
        },
        {
          abbr: "F. Name",
          title: "First Name",
        },
        {
          abbr: "L. Name",
          title: "Last Name"
        },
        {
          abbr: "Year",
          title: "Birth Year"
        }
      ],
      heatResultsRows: [],
    };
  },
  watch: {
    timerId: function () {
      this.heatFinish = []
      this.heatFinisher.clear()
      this.heatResultsRows = []
    }
  },
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants']),
    recordLaneFinish: function (lane) {
      if (!this.hasFinished(lane)) {
        this.heatFinish.push(lane)
        getParticipantByBib(this.$route.params.raceId, lane.bib).then((participant) => {
          if (Object.keys(participant).length !== 0)
            this.heatFinisher.set(lane.bib, participant)
        })
      }
    },
    bibInput: async function (e) {
      let ok = await recordResult(this.$route.params.raceId, this.$route.params.eventId, e.bib, this.timerId)

      if (ok) {
        if (e.success) {
          e.success()
        }
      } else {
        if (e.error) {
          e.error()
        }
      }
    },
    refreshData: function (e) {
      clearTimeout(this.resultsRefresh)

      listTimers(this.$route.params.raceId, this.$route.params.eventId).then((timers) => {
        this.timers = timers
        if (this.timerId && (this.timerId == "" || this.timerId == "latest")) {
          this.timerId = this.timers[0].id

          if (this.timers[0].assignments && this.timers[0].assignments.length > 0) {
            this.tabSelect("heat")
          }
        }

        if (this.timerId && this.timerId != "") {
          getHeatResults(this.$route.params.raceId, this.$route.params.eventId, this.timerId).then((results) => {
            this.results = results
            this.resultsRefresh = setTimeout(this.refreshData, 2500)
          })
        }
      })
    },
    formatMilliseconds,
    isActiveTab: function (tab) {
      return this.activeTab == tab;
    },
    tabSelect: function (tab) {
      this.activeTab = tab;
    },
    hasFinished: function (assignment) {
      for (const finisher of this.heatFinish) {
        if (finisher.lane == assignment.lane) {
          return true
        }
      }
      return false
    },
    saveHeat: async function () {
      for (const finisher of this.heatFinish) {
        let ok = await recordResult(this.$route.params.raceId, this.$route.params.eventId, finisher.bib, this.timerId)
        if (!ok) {
          return
        }
      }

      this.heatFinish = []
    }
  },
  computed: {
    ...mapState(useParticipantsStore, {
      first_name: (store) => (bib) => {
        const participant = store.participants.find((entry) => entry.bibNumber == bib)
        if (participant)
          return participant.firstName

        return '-'
      },
      last_name: (store) => (bib) => {
        const participant = store.participants.find((entry) => entry.bibNumber == bib)
        if (participant)
          return participant.lastName

        return '-'

      },
      birth_year: (store) => (bib) => {
        const participant = store.participants.find((entry) => entry.bibNumber == bib)
        if (participant)
          return participant.birthYear

        return '-'

      },
    }),
    heatResults: function () {
      let results = []
      let pos = 1
      for (const finishTime of this.results) {
        let finishedLane = this.heatFinish[pos - 1]
        let bib = finishTime.bib_number
        if (bib === "" && finishedLane) {
          bib = finishedLane.bib
        } else if (bib === "")
          bib = "-"

        if (!finishedLane)
          finishedLane = { bib: '-', lane: "-" }

        results.push([
          pos++,
          formatMilliseconds(finishTime.result_ms),
          finishedLane.lane,
          bib,
          this.first_name(bib),
          this.last_name(bib),
          this.birth_year(bib),
        ])
      }

      return results
    },
    hasLaneAssignments: function () {
      for (const timer of this.timers) {
        if (timer.id == this.timerId)
          return (timer.assignments && timer.assignments.length > 0)
      }
    },
    laneAssignments: function () {
      for (const timer of this.timers) {
        if (timer.id == this.timerId) {
          return timer.assignments
        }
      }

    },
    getFirstUnmatchedPlace: function () {
      for (let i = 0; i < this.results.length; i++) {
        let result = this.results[i]
        if (result.bib_number == "") {
          return i + 1
        }
      }

      return "✅"
    },
    getFirstUnmatchedTime: function () {
      for (let i = 0; i < this.results.length; i++) {
        let result = this.results[i]
        if (result.bib_number == "") {
          return formatMilliseconds(result.result_ms)
        }
      }

      return "None"
    },
    getHeatTotalResults: function () {
      return this.results.length
    }
  },
};
</script>

<style scoped>
@media only screen and (max-width: 768px) {
  .columns {
    display: flex;
    flex-direction: column;
    height: 80vh;
    /* Adjust as needed */
    overflow: hidden;
    /* Prevent page from scrolling */
  }

  .column.recorder {
    flex-grow: 0;
  }

  .column.results {
    flex: 1;
    /* Grow to fill remaining space */
    overflow-y: auto;
    overflow-x: auto;
    /* Enable vertical scrolling */
  }
}
</style>
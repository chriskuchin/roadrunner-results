<template>
  <div class="container">
    <h1>Lane Assignment</h1>
    <div class="section">
      <div class="field has-addons">
        <div class="control">
          <div class="select">
            <select name="country" v-model="id">
              <option value="">New Heat</option>
              <option v-for="(heat, index) in heats" :key="heat.timer_id" :value="heat.timer_id">
                Heat {{ index + 1 }}
              </option>
            </select>
          </div>
        </div>
        <div class="control">
          <button type="submit" class="button is-primary" @click="newHeat" v-if="id == ''">Create Heat</button>
          <button type="submit" class="button is-primary" @click="saveHeat" v-else>Update Heat</button>
        </div>
        <div class="control" v-if="id != ''">
          <button class="button is-danger" @click="deleteHeat">Delete</button>
        </div>
      </div>
      <table class="table is-narrow" style="margin: 0 auto;">
        <thead>
          <tr>
            <th><abbr title="Lane Assignment">Ln</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th v-if="view.includeAthleteFirstName"><abbr title="Athlete First Name">F. Name</abbr></th>
            <th v-if="view.includeAthleteLastName"><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Yr</abbr></th>
          </tr>
        </thead>
        <tfoot>
          <tr>
            <th><abbr title="Lane Assignment">Ln</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th v-if="view.includeAthleteFirstName"><abbr title="Athlete First Name">F. Name</abbr></th>
            <th v-if="view.includeAthleteLastName"><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Yr</abbr></th>
          </tr>
        </tfoot>
        <tbody>
          <tr v-for="( lane, index ) in lanes" :key="index" class="is-size-4">
            <th>{{ lane.lane }}</th>
            <td><input :tabindex="index + 1" class="input is-medium" type="number" placeholder="Bib Number"
                v-model="lanes[index].bib" style="max-width: 70px;"></td>
            <td v-if="view.includeAthleteFirstName" :class="{ 'truncate': view.shortFirstName }">
              {{ first_name(lane.bib) }}
            </td>
            <td v-if="view.includeAthleteLastName">
              {{ last_name(lane.bib) }}
            </td>
            <td>
              {{ birth_year(lane.bib) }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>

import { createNewHeat, listHeats, updateHeat, deleteHeat } from '../api/heats';
import { useParticipantsStore } from '../store/participants';
import { mapState, mapActions } from 'pinia';

export default {
  components: {},
  mounted: async function () {
    this.loadParticipants(this.$route.params.raceId, "", "", "", "", 500, 0)
    listHeats(this.$route.params.raceId, this.$route.params.eventId).then((heats) => {
      this.heats = heats
    })

    let vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0)

    if (vw < 380) {
      this.view.includeAthleteFirstName = false
      this.view.includeAthleteLastName = false
    } else if (vw < 425) {
      this.view.includeAthleteFirstName = true
      this.view.shortFirstName = true
      this.view.includeAthleteLastName = false
    } else if (vw > 500) {
      this.view.includeAthleteFirstName = true
      this.view.includeAthleteLastName = true
      this.view.shortFirstName = false
    }
  },
  data: function () {
    return {
      view: {
        includeAthleteLastName: false,
        includeAthleteFirstName: true,
        shortFirstName: true,
      },
      id: "",
      laneCount: 8,
      lanes: [],
      heats: [],
    }
  },
  watch: {
    id: {
      handler() {
        for (let heat of this.heats) {
          if (heat.timer_id == this.id) {
            this.activateHeat(heat)
            return
          }
        }

        this.activateHeat({
          timer_id: "",
          assignments: []
        })
      },
      immediate: true
    },
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
    })
  },
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants']),
    generateLaneAssignments(targetCount) {
      if (this.lanes.length < targetCount) {
        const lanesToAdd = targetCount - this.lanes.length
        const currentLength = this.lanes.length
        for (let i = 1; i <= lanesToAdd; i++) {
          this.lanes.push({
            lane: i + currentLength,
            bib: ""
          })
        }
      } else if (this.lanes.length > targetCount) {
        const lanesToRemove = this.lanes.length - targetCount
        for (let i = 1; i <= lanesToRemove; i++) {
          this.lanes.pop()
        }
      }
    },
    deleteHeat() {
      deleteHeat(this.$route.params.raceId, this.$route.params.eventId, this.id).then(() => {
        listHeats(this.$route.params.raceId, this.$route.params.eventId).then((heats) => {
          this.heats = heats
          this.id = ""
        })
      })
    },
    saveHeat() {
      updateHeat(this.$route.params.raceId, this.$route.params.eventId, this.id, this.lanes).then(() => {
        console.log("Finished")
      })
    },
    activateHeat(heat) {
      if (heat.assignments && heat.assignments.length !== 0)
        this.lanes = heat.assignments
      else
        this.lanes = []

      this.generateLaneAssignments(this.laneCount)
    },
    async newHeat() {
      const heatId = await createNewHeat(this.$route.params.raceId, this.$route.params.eventId, this.lanes)
      this.heats.push({
        assignments: this.lanes,
        timer_id: heatId
      })

      this.lanes = []
      this.generateLaneAssignments(this.laneCount)
    },
  }
};
</script>


<style scoped>
.truncate {
  white-space: nowrap;
  /* Prevents wrapping */
  overflow: hidden;
  /* Hides the overflowed content */
  text-overflow: ellipsis;
  /* Adds ellipsis (...) to indicate truncated text */
  max-width: 50px;
  /* Set maximum width to determine where to truncate */
}
</style>
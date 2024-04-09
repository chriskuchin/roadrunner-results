<template>
  <div class="container">
    <h1>Lane Assignment</h1>
    <div class="tabs">
      <ul>
        <li :class="{ 'is-active': heat.timer_id == id }" v-for="(heat, index) in heats" :key="heat.timer_id"
          @click="activateHeat(heat)">
          <a>
            Heat {{ index + 1 }}
          </a>
        </li>
        <li :class="{ 'is-active': id == '' }"><a @click="newHeat">New +</a></li>
      </ul>
    </div>
    <div class="section" v-if="id !== ''">
      <div class="select">
        <select v-model="laneCount">
          <option v-for=" n in 8 " :value="n + 3">{{ n + 3 }} Lanes</option>
        </select>
        <button class="button" @click="saveHeat">Save</button>
        <button class="button" @click="deleteHeat">Delete</button>
      </div>
      <table class="table" style="margin: 0 auto;">
        <thead>
          <tr>
            <th><abbr title="Lane Assignment">Lane</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th><abbr title="Athlete First Name">F. Name</abbr></th>
            <th><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Year</abbr></th>
          </tr>
        </thead>
        <tfoot>
          <tr>
            <th><abbr title="Lane Assignment">Lane</abbr></th>
            <th><abbr title="Athlete Bib Number">Bib</abbr></th>
            <th><abbr title="Athlete First Name">F. Name</abbr></th>
            <th><abbr title="Athlete Last Name">L. Name</abbr></th>
            <th><abbr title="Athlete Birth Year">B. Year</abbr></th>
          </tr>
        </tfoot>
        <tbody>
          <tr v-for="( lane, index ) in lanes " class="is-size-4">
            <th>{{ lane.lane }}</th>
            <td><input :tabindex="index + 1" class="input is-medium" type="text" placeholder="Bib Number"
                @blur="bibBlur" v-model="lanes[index].bib"></td>
            <td>{{ participants[index].first_name }}</td>
            <td>{{ participants[index].last_name }}</td>
            <td>{{ participants[index].birth_year }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>

import { createNewHeat, listHeats, updateHeat, deleteHeat } from '../api/heats';
import { getParticipantByBib } from '../api/participants';
import DropdownMenu from '../components/DropdownMenu.vue'

export default {
  components: {},
  mounted: async function () {
    this.heats = await listHeats(this.$route.params.raceId, this.$route.params.eventId)
  },
  data: function () {
    return {
      id: "",
      laneCount: 8,
      lanes: [],
      participants: [],
      heats: [],
    }
  },
  watch: {
    laneCount: {
      handler(newCount) {
        this.generateLaneAssignments(newCount)
      },
      immediate: true
    }
  },
  computed: {
  },
  methods: {
    generateLaneAssignments(targetCount) {
      if (this.lanes.length < targetCount) {
        const lanesToAdd = targetCount - this.lanes.length
        const currentLength = this.lanes.length
        for (let i = 1; i <= lanesToAdd; i++) {
          this.participants.push({
            first_name: "",
            last_name: "",
            birth_year: "",
          })
          this.lanes.push({
            lane: i + currentLength,
            bib: ""
          })
        }
      } else if (this.lanes.length > targetCount) {
        const lanesToRemove = this.lanes.length - targetCount
        for (let i = 1; i <= lanesToRemove; i++) {
          this.participants.pop()
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
      if (this.id != "") {
        console.log("update the heat")
      }
      this.id = heat.timer_id

      if (heat.assignments && heat.assignments.length !== 0)
        this.lanes = heat.assignments
      else
        this.lanes = []

      if (this.lanes.length == 0) {
        this.laneCount = 8
      } else {
        this.laneCount = this.lanes.length
      }

      this.generateLaneAssignments(this.laneCount)
    },
    async newHeat() {
      const heat = await createNewHeat(this.$route.params.raceId, this.$route.params.eventId, [])
      this.heats.push(heat)
    },
    async bibBlur(e) {
      const idx = e.target.tabIndex
      const bib = e.target.value

      if (bib !== "") {
        let participant = await getParticipantByBib(this.$route.params.raceId, bib)

        this.participants[idx - 1].first_name = participant.first_name
        this.participants[idx - 1].last_name = participant.last_name
        this.participants[idx - 1].birth_year = participant.birth_year
      }
    }
  }
};
</script>
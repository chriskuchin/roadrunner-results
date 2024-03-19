<template>
  <div class="container">
    <h1>Lane Assignment</h1>
    <div class="tabs">
      <ul>
        <li class="is-active"><a>Heat 1</a></li>
        <li><a>Heat 2</a></li>
        <li><a>Heat 3</a></li>
        <li><a>Heat 4</a></li>
        <li><a>New +</a></li>
      </ul>
    </div>
    <div class="section">
      <div class="select">
        <select v-model="laneCount">
          <option v-for="n in 8" :value="n + 3">{{ n + 3 }} Lanes</option>
        </select>
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
          <tr v-for="(lane, index) in lanes">
            <th>{{ lane.lane }}</th>
            <td><input :tabindex="index + 1" class="input is-small" type="text" placeholder="Bib Number" @blur="bibBlur"
                v-model="lanes[index].bib"></td>
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

import { getParticipantByBib } from '../api/participants';
import DropdownMenu from '../components/DropdownMenu.vue'

export default {
  mounted: function () {
    dm: DropdownMenu
  },
  data: function () {
    return {
      id: "123",
      laneCount: 8,
      lanes: [],
      participants: [],
    }
  },
  watch: {
    laneCount: {
      handler(newCount) {
        if (this.lanes.length < newCount) {
          const lanesToAdd = newCount - this.lanes.length
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
        } else if (this.lanes.length > newCount) {
          const lanesToRemove = this.lanes.length - newCount
          for (let i = 1; i <= lanesToRemove; i++) {
            this.participants.pop()
            this.lanes.pop()
          }
        }
      },
      immediate: true
    }
  },
  computed: {
  },
  methods: {
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
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
      <table class="table">
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
            <td><input :tabindex="index + 1" class="input is-small" type="text" placeholder="Text input"
                :value="lane.bib"></td>
            <td></td>
            <td></td>
            <td></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>

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
    }
  },
  watch: {
    laneCount: {
      handler(newCount) {
        if (this.lanes.length < newCount) {
          const lanesToAdd = newCount - this.lanes.length
          const currentLength = this.lanes.length
          for (let i = 1; i <= lanesToAdd; i++) {
            this.lanes.push({
              lane: i + currentLength,
              bib: ""
            })
          }
        } else if (this.lanes.length > newCount) {
          const lanesToRemove = this.lanes.length - newCount
          for (let i = 1; i <= lanesToRemove; i++) {
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
  }
};
</script>
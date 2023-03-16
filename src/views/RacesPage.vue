<template>
  <div>
    <div class="level">
      <div class="level-item">
        <table class="table">
          <thead>
            <tr>
              <th><abbr>Race</abbr></th>
              <th><abbr>Participants</abbr></th>
              <th><abbr>Events</abbr></th>
              <th><abbr>Results</abbr></th>
              <th><abbr>Delete</abbr></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="race in races" :key="race.id">
              <td><router-link :to="'/races/' + race.id + '/'">{{ race.name }}</router-link></td>
              <td><router-link :to="'/races/' + race.id + '/participants'">Participants</router-link></td>
              <td><router-link :to="'/races/' + race.id + '/events'">Events</router-link></td>
              <td><router-link :to="'/races/' + race.id + '/results'">Results</router-link></td>
              <td>
                <button class="button is-danger is-outlined is-small" @click="deleteRace(race.id)">
                  Delete
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="level">
      <div class="level-item">
        <div class="field has-addons">
          <div class="control">
            <input class="input" type="text" placeholder="Race Description" v-model="createRaceInput">
          </div>
          <div class="control">
            <button class="button is-info" :disabled="isCreateDisabled" @click="createRace">
              Create Race
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { getRaces, createRace, deleteRace } from '../api/races'

export default {
  data: function () {
    return {
      races: [],
      createRaceInput: ""
    };
  },
  mounted: function () {
    getRaces().then(res => this.races = res)
  },
  computed: {
    isCreateDisabled: function () {
      return this.createRaceInput == ""
    }
  },
  methods: {
    async deleteRace(raceID) {
      this.races = await deleteRace(raceID)
    },
    async createRace() {
      this.createRaceInput = ""
      this.races = await createRace(this.createRaceInput)
    },
  },
};
</script>
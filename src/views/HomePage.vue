<template>
  <div class="section">
    <div class="box" v-for="race in racesStore.getRaces" :key="race.id">
      <div class="title is-4">{{ race.name }}</div>
      <router-link :to="'/races/' + race.id + '/events'" class="button">Events</router-link>
      <!-- <button class="delete" @click="deleteRace(race.id)"></button> -->
    </div>
    <fab @click="toggleCreateRaceModal"></fab>
    <modal @close="toggleCreateRaceModal" :show="raceModal.show">
      <p class="title">Create Race</p>
      <div class="field">
        <label class="label">Description</label>
        <div class="control">
          <input class="input" type="text" placeholder="Race Description" v-model="raceModal.description">
        </div>
      </div>
      <div class="field is-grouped">
        <div class="control">
          <button :class="['button', 'is-link', { 'is-loading': raceModal.creating }]" @click="createRace">Submit</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="toggleCreateRaceModal">Cancel</button>
        </div>
      </div>
    </modal>
  </div>
</template>

<script>
import { mapStores } from 'pinia'
import { useRacesStore } from "../store/races";
import { createRace, deleteRace } from "../api/races"
import Modal from '../components/Modal.vue'
import FAB from '../components/Fab.vue'

export default {
  components: {
    "modal": Modal,
    "fab": FAB,
  },
  data: function () {
    return {
      raceModal: {
        show: false,
        creating: false,
        description: "",
      }
    }
  },
  methods: {
    createRace: function () {
      var self = this
      self.raceModal.creating = true
      createRace(self.raceModal.description).then(() => {
        self.raceModal.creating = false
        self.raceModal.description = ""
        self.toggleCreateRaceModal()
        self.racesStore.loadRaces()
      })

    },
    cancelCreateRace: function () {
      this.raceModal.description = ""
    },
    deleteRace: function (raceID) {
      var self = this
      deleteRace(raceID).then(() => {
        self.racesStore.loadRaces()
      })
    },
    toggleCreateRaceModal: function () {
      this.raceModal.show = !this.raceModal.show

    }
  },
  computed: {
    ...mapStores(useRacesStore),
    raceLink: function (id) {
      return "/races/" + id + "/results"
    }
  },
  mounted: function () {
    this.racesStore.loadRaces()
  }
};
</script>
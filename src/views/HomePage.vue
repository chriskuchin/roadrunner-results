<template>
  <div class="section">
    <div class="box mb-4" v-for="race in racesStore.races" :key="race.id">
      <div class="has-text-right">
        <cm class="is-right" v-if="isLoggedIn">
          <!-- <a class="dropdown-item">Edit Race</a> -->
          <router-link class="dropdown-item" :to="'/races/' + race.id + '/volunteers'">Volunteers</router-link>
          <a class="dropdown-item" @click="openShareModal(race.id)">Add Volunteer</a>
          <a class="dropdown-item" @click="openDivisionsModal(race.id)">Manage Divisions</a>
          <a class="dropdown-item" @click="generateDivisions(race.id)">Generate Divisions</a>
          <router-link class="dropdown-item" :to="'/races/' + race.id">Race Info</router-link>
          <hr class="dropdown-divider" />
          <a class="dropdown-item" @click="deleteRace(race.id)">Delete Race</a>
        </cm>
      </div>

      <div class="title is-4">{{ race.name }}</div>
      <div class="subtitle mb-4" v-if="race.date">{{ getRaceDate(race.date) }}</div>
      <div class="buttons">
        <router-link :to="'/races/' + race.id + '/events'"
          class="button is-info is-light is-outlined">Events</router-link>
        <router-link :to="'/races/' + race.id + '/participants'"
          class="button is-info is-light is-outlined">Participants</router-link>
      </div>
    </div>
    <fab v-if="isLoggedIn" @click="openCreateRaceModal"></fab>
    <create-race-modal ref="create-race-modal" />
    <share-modal ref="share-modal" />
    <divisions-modal ref="divisions-modal" />
    <not :show="error.show" type="is-danger is-light" @close="dismissError">{{ error.msg }}</not>
  </div>
</template>

<script>
import { mapStores, mapState, mapActions } from 'pinia'
import { useRacesStore } from "../store/races"
import { useRaceStore } from '../store/race'
import { useUserStore } from '../store/user'
import { useDivisionsStore } from '../store/divisions'
import { deleteRace } from "../api/races"

import FAB from '../components/Fab.vue'
import Notification from '../components/Notification.vue'
import ContextMenu from '../components/DropdownMenu.vue'
import ShareRaceModal from '../components/modals/ShareModal.vue'
import DivisionsModal from '../components/modals/DivisionsModal.vue'
import CreateRaceModal from '../components/modals/CreateRaceModal.vue'

export default {
  components: {
    "share-modal": ShareRaceModal,
    "divisions-modal": DivisionsModal,
    "create-race-modal": CreateRaceModal,
    "fab": FAB,
    "not": Notification,
    "cm": ContextMenu,
  },
  data: function () {
    return {
      error: {
        show: false,
        msg: "",
      },
    }
  },
  methods: {
    getRaceMenu: function (e) {
      let target = e.currentTarget
      if (target.className.includes("is-active")) {
        target.className = target.className.replace("is-active", "")
      }
      else {
        target.className = e.currentTarget.className + " is-active"
      }
    },
    getRaceDate: function (date) {
      const options = {
        weekday: 'long',
        year: 'numeric',
        month: 'long',
        day: 'numeric',
      };

      let dateParts = date.split("-", 3)
      let raceDate = new Date(dateParts[0], dateParts[1] - 1, dateParts[2])
      return raceDate.toLocaleDateString('us-EN', options)
    },
    contextMenuClick: function (raceId) {
      if (this.eventContextMenu.active[raceId] === undefined)
        this.eventContextMenu.active[raceId] = false

      this.eventContextMenu.active[raceId] = !this.eventContextMenu.active[raceId]
    },
    dismissError: function () {
      this.error.show = false
    },
    deleteRace: function (raceID) {
      var self = this
      deleteRace(raceID).then(() => {
        self.racesStore.loadRaces()
      }).catch((err) => {
        self.error.show = true
        self.error.msg = err
      })
    },
    openCreateRaceModal: function () {
      this.$refs['create-race-modal'].open()
    },
    openShareModal: function (raceId) {
      this.$refs['share-modal'].open(raceId)
    },
    openDivisionsModal: function (raceId) {
      this.$refs['divisions-modal'].open(raceId)
    },
    generateDivisions: async function (raceID) {
      let year = new Date().getFullYear()
      let firstDivision = `${year - 6}+`
      let yougestDivisionFilter = []
      for (var i = year - 6; i < year; i++) {
        yougestDivisionFilter.push(`${i}`)
      }
      let divisions = {}
      divisions[firstDivision] = yougestDivisionFilter

      for (var i = 1; i <= 30; i = i + 2) {
        var high = year - (6 + i)
        var low = year - (6 + i + 1)
        var currentDivision = `${low}-${high}`

        divisions[currentDivision] = [`${low}`, `${high}`]
      }
      let genders = ["Male", "Female"]
      Object.keys(divisions).forEach((desc) => {
        genders.forEach((gender) => {
          this.createDivision(raceID, `${desc} ${gender}`, [gender], divisions[desc])
        })
      })
    },
    ...mapActions(useDivisionsStore, ['createDivision', 'load']),
    ...mapActions(useRaceStore, ['shareRace', 'listVolunteers'])
  },
  computed: {
    ...mapStores(useRacesStore),
    ...mapState(useUserStore, ['isLoggedIn']),
    ...mapState(useDivisionsStore, ['divisions']),
    raceLink: function (id) {
      return "/races/" + id + "/results"
    }
  },
  mounted: function () {
    this.racesStore.loadRaces()
  }
};
</script>
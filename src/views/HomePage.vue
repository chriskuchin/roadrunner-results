<template>
  <div class="section">
    <div class="box" v-for="race in racesStore.getRaces" :key="race.id">
      <div class="has-text-right">
        <div class="dropdown is-right" @click="getRaceMenu" v-if="isLoggedIn">
          <div class="dropdown-trigger">
            <span class="icon is-clickable" aria-haspopup="true" aria-controls="dropdown-menu">
              <icon icon="fa-solid fa-ellipsis-v"></icon>
            </span>
          </div>
          <div class="dropdown-menu" id="dropdown-menu" role="menu">
            <div class="dropdown-content">
              <a class="dropdown-item">Edit Race</a>
              <router-link class="dropdown-item" :to="'/races/' + race.id + '/volunteers'">Volunteers</router-link>
              <a class="dropdown-item" @click="toggleShareModal(race.id)">Add Volunteer</a>
              <a class="dropdown-item" @click="toggleDivisionModal(race.id)">Manage Divisions</a>
              <a class="dropdown-item" @click="generateDivisions(race.id)">Generate Divisions</a>
              <router-link class="dropdown-item" :to="'/races/' + race.id">Race Info</router-link>
              <hr class="dropdown-divider" />
              <a class="dropdown-item" @click="deleteRace(race.id)">Delete Race</a>
            </div>
          </div>
        </div>
      </div>

      <div class="title is-4">{{ race.name }}</div>
      <div class="subtitle" v-if="race.date">{{ getRaceDate(race.date) }}</div>
      <div class="buttons">
        <router-link :to="'/races/' + race.id + '/events'"
          class="button is-info is-light is-outlined">Events</router-link>
        <router-link :to="'/races/' + race.id + '/participants'"
          class="button is-info is-light is-outlined">Participants</router-link>
      </div>
    </div>
    <fab v-if="isLoggedIn" @click="toggleCreateRaceModal"></fab>
    <modal @close="toggleShareModal" :show="volunteerModal.show">
      <div class="title">Volunteers</div>
      <hr />
      <div class="subtitle">Add Volunteers</div>
      <span class="tag is-black mr-2" v-for="volunteer in volunteerModal.volunteers" :key="volunteer.userId">
        {{ volunteer.email }}
      </span>
      <div class="field">
        <label class="label">Label</label>
        <div class="control">
          <input class="input" type="text" placeholder="Text input" v-model="volunteerModal.emailInput"
            v-on:keyup.enter="appendEmailForShare">
        </div>
        <p class="help">Volunteers must have a user account. Enter the users email address above.</p>
      </div>
      <div class="field is-grouped is-grouped-multiline">
        <div class="tags has-addons" v-for="email in volunteerModal.emails" style="margin: 0 5px;">
          <span class="tag is-link">{{ email.email }}</span>
          <span class="tag is-delete" @click="removeVolunteerEmail(email.email)"></span>
        </div>
      </div>
      <div class="field is-grouped">
        <div class="control">
          <button class="button is-link" @click="addVolunteers">Share</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="toggleShareModal">Cancel</button>
        </div>
      </div>
    </modal>
    <modal @close="toggleCreateRaceModal" :show="raceModal.show">
      <div class="tabs">
        <ul>
          <li :class="{ 'is-active': raceModal.tabs.activeTab == 'create' }"><a
              @click="raceModalTabClick('create')">Create Race</a></li>
          <li :class="{ 'is-active': raceModal.tabs.activeTab == 'import' }"><a
              @click="raceModalTabClick('import')">Import Race</a></li>
        </ul>
      </div>
      <div v-if="raceModal.tabs.activeTab == 'create'">
        <p class="title">Create Race</p>
        <div class="field">
          <label class="label">Description</label>
          <div class="control">
            <input class="input" type="text" placeholder="Race Description" v-model="raceModal.description">
          </div>
        </div>
        <div class="field">
          <label class="label">Date</label>
          <div class="control">
            <input class="input" type="date" placeholder="" v-model="raceModal.date"> {{ raceModal.date }}
          </div>
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button :class="['button', 'is-link', { 'is-loading': raceModal.creating }]"
              @click="createRace">Submit</button>
          </div>
          <div class="control">
            <button class="button is-link is-light" @click="toggleCreateRaceModal">Cancel</button>
          </div>
        </div>
      </div>
      <div v-if="raceModal.tabs.activeTab == 'import'">
        <p class="title">Import Race</p>
        <div class="field">
          <label class="label">Race URL</label>
          <div class="control">
            <input class="input" type="text" placeholder="Race URL" v-model="raceModal.importURL">
          </div>
        </div>
        <div class="field">
          <label class="label">Description</label>
          <div class="control">
            <input class="input" type="text" placeholder="Race Description" v-model="raceModal.description">
          </div>
        </div>
        <div class="field">
          <label class="label">Date</label>
          <div class="control">
            <input class="input" type="date" placeholder="" v-model="raceModal.date">
          </div>
        </div>
        <div class="field is-grouped">
          <div class="control">
            <button :class="['button', 'is-link', { 'is-loading': raceModal.importing }]"
              @click="raceModalImportRace">Submit</button>
          </div>
          <div class="control">
            <button class="button is-link is-light" @click="toggleCreateRaceModal">Cancel</button>
          </div>
        </div>
      </div>
    </modal>
    <modal @close="toggleDivisionModal" :show="divisionModal.show">
      <div v-for="division in divisions">{{ division.display }} - {{ division.filters }}</div>
    </modal>
    <not :show="error.show" type="is-danger is-light" @close="dismissError">{{ error.msg }}</not>
  </div>
</template>

<script>
import { mapStores, mapState, mapActions } from 'pinia'
import { useRacesStore } from "../store/races"
import { useRaceStore } from '../store/race'
import { useUserStore } from '../store/user'
import { useDivisionsStore } from '../store/divisions'
import { createRace, importRace, deleteRace } from "../api/races"
import Modal from '../components/Modal.vue'
import FAB from '../components/Fab.vue'
import Notification from '../components/Notification.vue'

export default {
  components: {
    "modal": Modal,
    "fab": FAB,
    "not": Notification,
  },
  data: function () {
    return {
      divisionModal: {
        show: false,
        raceId: "",
      },
      volunteerModal: {
        show: false,
        emailInput: "",
        raceID: "",
        emails: [],
        volunteers: [],
      },
      volunteersModal: {
        show: false,
      },
      raceModal: {
        show: false,
        creating: false,
        importing: false,
        description: "",
        date: "",
        importURL: "",
        tabs: {
          activeTab: "create",
        }
      },
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
    raceModalImportRace: async function () {
      this.raceModal.importing = true
      await importRace(this.raceModal.importURL, this.raceModal.description, this.raceModal.date)
      this.toggleCreateRaceModal()
      this.racesStore.loadRaces()
      this.raceModal.importing = false
    },
    raceModalTabClick: function (tabID) {
      this.raceModal.tabs.activeTab = tabID
    },
    contextMenuClick: function (raceId) {
      if (this.eventContextMenu.active[raceId] === undefined)
        this.eventContextMenu.active[raceId] = false

      this.eventContextMenu.active[raceId] = !this.eventContextMenu.active[raceId]
    },
    dismissError: function () {
      this.error.show = false
    },
    createRace: function () {
      var self = this
      self.raceModal.creating = true
      createRace(self.raceModal.description, self.raceModal.date).then(() => {
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
      }).catch((err) => {
        self.error.show = true
        self.error.msg = err
      })
    },
    toggleDivisionModal: function (raceID) {
      this.divisionModal.show = !this.divisionModal.show
      this.load(raceID).then(() => {
        console.log("Loaded Divisions")
      })
    },
    toggleCreateRaceModal: function () {
      this.raceModal.show = !this.raceModal.show
    },
    toggleShareModal: function (raceID) {
      this.volunteerModal.show = !this.volunteerModal.show

      if (raceID && (typeof raceID === 'string' || raceID instanceof String)) {
        this.volunteerModal.raceID = raceID
      }

      if (!this.volunteerModal.show) {
        this.volunteerModal.emails = []
        this.volunteerModal.emailInput = ""
        this.volunteerModal.raceID = ""
      } else {
        this.listVolunteers(raceID).then((volunteers) => {
          this.volunteerModal.volunteers = volunteers
        })
      }
    },
    appendEmailForShare: function () {
      let email = this.volunteerModal.emailInput.trim()
      if (email == "")
        return

      this.volunteerModal.emailInput = ""
      this.volunteerModal.emails.push({
        email: email
      })
    },
    removeVolunteerEmail: function (email) {
      let loc = this.volunteerModal.emails.indexOf(email)
      if (loc > -1) {
        this.volunteerModal.emails.splice(loc, 1)
      }
    },
    addVolunteers: async function () {
      await this.shareRace(this.volunteerModal.raceID, this.volunteerModal.emails)

      this.toggleShareModal()
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

      for (var i = 1; i <= 10; i = i + 2) {
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
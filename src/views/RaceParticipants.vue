<template>
  <div class="container mx-5">
    <h1 class="title">Participants</h1>
    <div class="table-container">
      <table class="table" style="margin: 0 auto;">
        <thead>
          <tr>
            <th><abbr title="Bib Number">Bib</abbr></th>
            <th><abbr title="First & Last Name">Name</abbr></th>
            <th><abbr title="Birth Year">Year</abbr></th>
            <th>Gender</th>
            <th>Team</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="participant in currentParticipants">
            <td>{{ participant.bibNumber }}</td>
            <td>{{ participant.firstName }} {{ participant.lastName }}</td>
            <td>{{ participant.birthYear }}</td>
            <td>{{ participant.gender }}</td>
            <td>{{ participant.team }}</td>
            <td><button class="button is-small is-light is-info" @click="editParticipant(participant)">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <fab @click="fabAction" v-if="isLoggedIn" />

    <modal :show="editModal.show" @close="toggleEditModal">
      <h1 class="title">Edit Participant</h1>
      <part-form :first_name="editForm.first_name" :last_name="editForm.last_name" :bib_number="editForm.bib_number"
        :birth_year="editForm.birth_year" :gender="editForm.gender" :team="editForm.team"
        @update:bib_number="editForm.bib_number = Number($event)"
        @update:birth_year="editForm.birth_year = Number($event)" @update:first_name="editForm.first_name = $event"
        @update:gender="editForm.gender = $event" @update:last_name="editForm.last_name = $event"
        @update:team="editForm.team = $event" />
      <div class="field is-horizontal mt-5 mb-0 is-pulled-right">
        <div class=" field-body">
          <div class="field is-grouped">
            <div class="control">
              <button class="button is-link" @click="submitEditModal">Submit</button>
            </div>
            <div class="control">
              <button class="button is-link is-light" @click="toggleEditModal">Cancel</button>
            </div>
          </div>
        </div>
      </div>
      <div class="is-clearfix"></div>
    </modal>

    <modal :show="registerModal.show" @close="toggleRegisterModal">
      <h1 class="title">Register</h1>
      <div class="tabs">
        <ul>
          <li :class="{ 'is-active': activeTab('register') }"><a @click="clickTab('register')">Register Participant</a>
          </li>
          <li :class="{ 'is-active': activeTab('upload') }"> <a @click=" clickTab('upload')">Upload CSV</a></li>
        </ul>
      </div>
      <div class="container upload-csv" v-if="registerModal.tabs.active == 'upload'">
        <div class="field is-horizontal">
          <div class="file has-name">
            <label class="file-label">
              <input class="file-input" type="file" name="resume" @change="selectFile">
              <span class="file-cta">
                <span class="file-icon">
                  <i class="fas fa-upload"></i>
                </span>
                <span class="file-label">
                  Choose a file…
                </span>
              </span>
              <span class="file-name">
                {{ uploadCSVForm.fileName }}
              </span>
            </label>
          </div>
        </div>
      </div>
      <div class="container register-participant" v-else-if="registerModal.tabs.active == 'register'">
        <part-form :first_name="registerForm.first_name" :last_name="registerForm.last_name"
          :bib_number="registerForm.bib_number" :birth_year="registerForm.birth_year" :gender="registerForm.gender"
          :team="registerForm.team" @update:bib_number="registerForm.bib_number = $event"
          @update:birth_year="registerForm.birth_year = $event" @update:first_name="registerForm.first_name = $event"
          @update:gender="registerForm.gender = $event" @update:last_name="registerForm.last_name = $event"
          @update:team="registerForm.team = $event" />
      </div>
      <div class="field is-horizontal mt-5 mb-0 is-pulled-right">
        <div class=" field-body">
          <div class="field is-grouped">
            <div class="control">
              <button class="button is-link" @click="submitRegisterModal">Submit</button>
            </div>
            <div class="control">
              <button class="button is-link is-light" @click="toggleRegisterModal">Cancel</button>
            </div>
          </div>
        </div>
      </div>
      <div class="is-clearfix"></div>
    </modal>
  </div>
</template>

<script>
import FAB from '../components/Fab.vue'
import Modal from '../components/Modal.vue'
import ParticipantForm from '../components/ParticipantForm.vue'

import { mapActions, mapState } from 'pinia'
import { useUserStore } from '../store/user'
import { useParticipantsStore } from '../store/participants'

export default {
  components: {
    'fab': FAB,
    'modal': Modal,
    'part-form': ParticipantForm,
  },
  data: function () {
    return {
      filters: {
        offset: 0,
        limit: 500,
        name: "",
        team: "",
        gender: "",
        year: "",
      },
      registerModal: {
        show: false,
        tabs: {
          active: "register"
        }
      },
      editModal: {
        show: false,
      },
      editForm: {
        id: "",
        first_name: "",
        last_name: "",
        bib_number: "",
        birth_year: "",
        team: "",
        gender: "",
      },
      uploadCSVForm: {
        file: null,
        fileName: "Please Select a File...",
        defaultFileName: "Please Select a File..."
      },
      registerForm: {
        first_name: "",
        last_name: "",
        bib_number: "",
        birth_year: new Date().getFullYear(),
        team: "",
        gender: "Female",
      },
    }
  },
  mounted: function () {
    this.loadParticipants(this.$route.params.raceId, this.filters.gender, this.filters.team, this.filters.team, this.filters.name, this.filters.limit, this.filters.offset)
  },
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants', 'registerParticipant', 'uploadParticipantCSV', 'updateParticipant']),
    ...mapActions(useUserStore, ['isLoggedIn']),
    fabAction: function () {
      this.toggleRegisterModal()
    },
    toggleRegisterModal: function () {
      this.registerModal.show = !this.registerModal.show

      if (!this.registerModal.show) {
        this.resetRegisterForm()
      }
    },
    submitRegisterModal: async function () {
      var success
      if (this.activeTab('register')) {
        success = await this.registerParticipant(this.$route.params.raceId, this.registerForm.first_name, this.registerForm.last_name, this.registerForm.bib_number, this.registerForm.gender, this.registerForm.birth_year, this.registerForm.team)
      } else if (this.activeTab('upload')) {
        success = await this.uploadParticipantCSV(this.$route.params.raceId, this.uploadCSVForm.file)
      }

      if (success) {
        this.reload()
        this.toggleRegisterModal()
        this.resetRegisterForm()
      } else {
        // Need to introduce the Error Store/Handler
        console.log("Failed to register participant")
      }
    },
    cancelRegistration: function () {
      this.toggleRegisterModal()
      this.resetRegisterForm()
    },
    reload: function () {
      this.loadParticipants(this.$route.params.raceId, this.filters.gender, this.filters.team, this.filters.team, this.filters.name, this.filters.limit, this.filters.offset)
    },
    resetRegisterForm: function () {
      this.registerModal.tabs.active = 'register'

      this.registerForm.first_name = ""
      this.registerForm.last_name = ""
      this.registerForm.bib_number = ""
      this.registerForm.birth_year = new Date().getFullYear()
      this.registerForm.team = ""
      this.registerForm.gender = "Female"

      this.uploadCSVForm.fileName = this.uploadCSVForm.defaultFileName
      this.uploadCSVForm.file = null
    },
    selectFile: function (e) {
      this.uploadCSVForm.file = e.target.files[0]
      this.uploadCSVForm.fileName = e.target.files[0].name
    },
    clickTab: function (tab) {
      this.registerModal.tabs.active = tab
    },
    activeTab: function (tab) {
      return tab === this.registerModal.tabs.active
    },
    toggleEditModal: function () {
      this.editModal.show = !this.editModal.show

      if (!this.editModal.show) {
        this.editForm.id = ""
        this.editForm.first_name = ""
        this.editForm.last_name = ""
        this.editForm.bib_number = ""
        this.editForm.gender = ""
        this.editForm.team = ""
        this.editForm.birth_year = ""
      }
    },
    editParticipant: function (participant) {
      this.editForm.first_name = participant.firstName
      this.editForm.last_name = participant.lastName
      this.editForm.bib_number = participant.bibNumber
      this.editForm.gender = participant.gender
      this.editForm.team = participant.team
      this.editForm.birth_year = participant.birthYear
      this.editForm.id = participant.id
      this.toggleEditModal()
    },
    submitEditModal: async function () {
      let success = await this.updateParticipant(this.$route.params.raceId, this.editForm.id, this.editForm)

      if (success) {
        this.toggleEditModal()
        this.reload()
      } else {
        console.log("fail")
      }
    }
  },
  computed: {
    ...mapState(useParticipantsStore, ['currentParticipants']),
    raceID: function () {
      return this.$route.params.raceId
    },
  }
}
</script>
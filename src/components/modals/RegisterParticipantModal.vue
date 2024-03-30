<template>
  <modal ref="modal" @close="reset">
    <h1 class="title mb-4">Register</h1>
    <div class="tabs">
      <ul>
        <li :class="{ 'is-active': activeTab('register') }"><a @click="clickTab('register')">Register Participant</a>
        </li>
        <li :class="{ 'is-active': activeTab('upload') }"> <a @click=" clickTab('upload')">Upload CSV</a></li>
      </ul>
    </div>
    <div class="container upload-csv" v-if="registerModal.tabs.active == 'upload'">
      <div class="field is-horizontal mt-4">
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
    <div class="field is-horizontal mt-5 is-pulled-right">
      <div class=" field-body">
        <div class="field is-grouped">
          <div class="control">
            <button class="button is-link" @click="submit">Submit</button>
          </div>
          <div class="control">
            <button class="button is-link is-light" @click="cancel">Cancel</button>
          </div>
        </div>
      </div>
    </div>
    <div class="is-clearfix"></div>
  </modal>
</template>

<script>
import { mapActions } from 'pinia'
import { useParticipantsStore } from '../../store/participants'

import ParticipantForm from '../ParticipantForm.vue'
import Modal from '../Modal.vue'

export default {
  components: {
    modal: Modal,
    'part-form': ParticipantForm
  },
  emits: ['reload'],
  data: function () {
    return {
      registerModal: {
        show: false,
        tabs: {
          active: "register"
        }
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
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants', 'registerParticipant', 'uploadParticipantCSV', 'updateParticipant']),
    submit: async function () {
      var success
      if (this.activeTab('register')) {
        success = await this.registerParticipant(this.$route.params.raceId, this.registerForm.first_name, this.registerForm.last_name, this.registerForm.bib_number, this.registerForm.gender, this.registerForm.birth_year, this.registerForm.team)
      } else if (this.activeTab('upload')) {
        success = await this.uploadParticipantCSV(this.$route.params.raceId, this.uploadCSVForm.file)
      }

      if (success) {
        this.reload()
        this.toggle()
        this.reset()
      } else {
        // Need to introduce the Error Store/Handler
        console.log("Failed to register participant")
      }
    },
    reload: function () {
      this.$emit('reload')
    },
    cancel: function () {
      this.toggle()
      this.reset()
    },
    toggle: function () {
      this.$refs.modal.toggle()
    },
    reset: function () {
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
  }
}
</script>
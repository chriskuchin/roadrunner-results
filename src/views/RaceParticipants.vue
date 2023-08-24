<template>
  <div class="container">
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
          </tr>
        </thead>
        <tbody>
          <tr v-for="participant in currentParticipants">
            <td>{{ participant.bibNumber }}</td>
            <td>{{ participant.firstName }} {{ participant.lastName }}</td>
            <td>{{ participant.birthYear }}</td>
            <td>{{ participant.gender }}</td>
            <td>{{ participant.team }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <fab @click="fabAction" v-if="isLoggedIn" />

    <modal :show="registerModal.show" @close="toggleRegisterModal">
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Name</label>
        </div>
        <div class="field-body">
          <div class="field">
            <p class="control is-expanded">
              <input class="input" type="text" placeholder="First Name" v-model="registerForm.first_name">
            </p>
          </div>
          <div class="field">
            <p class="control is-expanded">
              <input class="input" type="tex" placeholder="Last Name" v-model="registerForm.last_name">
            </p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Race Info</label>
        </div>
        <div class="field-body">
          <div class="field">
            <p class="control is-expanded">
              <input class="input" type="text" placeholder="Bib Number" v-model="registerForm.bib_number">
            </p>
          </div>
          <div class="field">
            <p class="control is-expanded">
              <input class="input" type="text" placeholder="Team" v-model="registerForm.team">
            </p>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Gender</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <label class="radio">
                <input type="radio" name="gender" value="Male" v-model="registerForm.gender">
                Male
              </label>
              <label class="radio">
                <input type="radio" name="gender" value="Female" v-model="registerForm.gender">
                Female
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">Birth Year</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control is-expanded">
              <div class="select">
                <select v-model="registerForm.birth_year">
                  <option>{{ year }}</option>
                  <option v-for="n in 80">{{ year - n }}</option>
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>


      <div class="field is-horizontal">
        <div class="field-label">
          <!-- Left empty for spacing -->
        </div>
        <div class="field-body">
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
    </modal>
  </div>
</template>

<script>
import FAB from '../components/Fab.vue'
import Modal from '../components/Modal.vue'

import { mapActions, mapState } from 'pinia'
import { useUserStore } from '../store/user'
import { useParticipantsStore } from '../store/participants'

export default {
  components: {
    'fab': FAB,
    'modal': Modal,
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
    ...mapActions(useParticipantsStore, ['loadParticipants', 'registerParticipant']),
    ...mapActions(useUserStore, ['isLoggedIn']),
    fabAction: function () {
      this.toggleRegisterModal()
    },
    toggleRegisterModal: function () {
      this.registerModal.show = !this.registerModal.show
    },
    submitRegisterModal: async function () {
      let success = await this.registerParticipant(this.$route.params.raceId, this.registerForm.first_name, this.registerForm.last_name, this.registerForm.bib_number, this.registerForm.gender, this.registerForm.birth_year, this.registerForm.team)
      if (success) {
        this.loadParticipants(this.$route.params.raceId, this.filters.gender, this.filters.team, this.filters.team, this.filters.name, this.filters.limit, this.filters.offset)
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
    resetRegisterForm: function () {
      this.registerForm.first_name = ""
      this.registerForm.last_name = ""
      this.registerForm.bib_number = ""
      this.registerForm.birth_year = new Date().getFullYear()
      this.registerForm.team = ""
      this.registerForm.gender = "Female"
    }
  },
  computed: {
    ...mapState(useParticipantsStore, ['currentParticipants']),
    year: () => new Date().getFullYear(),
    raceID: function () {
      return this.$route.params.raceId
    },
  }
}
</script>
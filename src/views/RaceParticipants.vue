<template>
  <div class="section">
    <h1 class="title">Participants</h1>
    {{ raceID }}
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
                <input type="radio" name="gender" v-model="registerForm.gender">
                Male
              </label>
              <label class="radio">
                <input type="radio" name="gender" v-model="registerForm.gender">
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
              <button class="button is-link" @click="createParticipant">Submit</button>
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

import { useUserStore } from '../store/user'
import { mapActions } from 'pinia'

export default {
  components: {
    'fab': FAB,
    'modal': Modal,
  },
  data: function () {
    return {
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
  methods: {
    ...mapActions(useUserStore, ['isLoggedIn']),
    fabAction: function () {
      this.toggleRegisterModal()
    },
    toggleRegisterModal: function () {
      this.registerModal.show = !this.registerModal.show
    },
    registerParticipant: function () {

    },
    cancelRegistration: function () {
      this.toggleRegisterModal()

      // clear form
    }
  },
  computed: {
    year: () => new Date().getFullYear(),
    raceID: function () {
      return this.$route.params.raceId
    },
  }
}
</script>
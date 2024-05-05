<template>
  <modal ref="share-modal" @open="opened">
    <div class="title">Add Volunteers</div>
    <span class="tag is-black mr-2" v-for="volunteer in volunteers" :key="volunteer.userId">
      {{ volunteer.email }}
    </span>
    <div class="field">
      <label class="label">Label</label>
      <div class="control">
        <input class="input" type="text" placeholder="Text input" v-model="emailInput"
          v-on:keyup.enter="appendEmailForShare">
      </div>
      <p class="help">Volunteers must have a user account. Enter the users email address above.</p>
    </div>
    <div class="field is-grouped is-grouped-multiline">
      <div class="tags has-addons" v-for="(email, idx) in emails" style="margin: 0 5px;">
        <span class="tag is-link">{{ email.email }}</span>
        <span class="tag is-delete" @click="removeVolunteerEmail(idx)"></span>
      </div>
    </div>
    <div class="field is-grouped mt-4">
      <div class="control">
        <button class="button is-link" @click="addVolunteers">Share</button>
      </div>
      <div class="control">
        <button class="button is-link is-light" @click="close">Cancel</button>
      </div>
    </div>
  </modal>
</template>

<script>
import Modal from '../Modal.vue'
import { useErrorBus } from '../../store/error'
import { useRaceStore } from '../../store/race'
import { mapActions } from 'pinia'
import { addRaceVolunteers } from '../../api/races'

export default {
  components: {
    modal: Modal,
  },
  data: function () {
    return {
      emailInput: "",
      raceID: "",
      emails: [],
      volunteers: []
    }
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    ...mapActions(useRaceStore, ['loadVolunteers']),
    open: function (raceID) {
      this.$refs['share-modal'].toggle()

      if (raceID && (typeof raceID === 'string' || raceID instanceof String)) {
        this.raceID = raceID
      }
    },
    close: function () {
      this.$refs['share-modal'].toggle()

      this.clear()
    },
    submit: function () {
      this.listVolunteers(this.raceID).then((volunteers) => {
        this.volunteers = volunteers
      })
    },
    clear: function () {
      this.emails = []
      this.emailInput = ""
      this.raceID = ""

    },
    opened: function () {
      this.clear()
    },
    appendEmailForShare: function () {
      let email = this.emailInput.trim()
      if (email == "")
        return

      this.emailInput = ""
      this.emails.push({
        email: email
      })
    },
    removeVolunteerEmail: function (loc) {
      if (loc > -1) {
        this.emails.splice(loc, 1)
      }
    },
    addVolunteers: async function () {
      addRaceVolunteers(this.raceID, this.emails).then((result) => {
        if (result.failure.length > 0) {
          this.handleError(`Failed to add volunteers: ${result.failure.join(", ")}`)
        }

        if (result.success.length > 0) {
          this.loadVolunteers()
        }

        this.$refs['share-modal'].toggle()
      }).catch((err) => {
        this.handleError(err)
      })

    },
  }
}
</script>
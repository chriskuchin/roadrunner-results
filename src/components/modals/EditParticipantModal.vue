<template>
  <modal ref="modal" @close="reset">
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
            <button class="button is-link" @click="submit">Submit</button>
          </div>
          <div class="control">
            <button class="button is-link is-light" @click="toggle">Cancel</button>
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

import Modal from '../Modal.vue';
import ParticipantForm from '../ParticipantForm.vue'

export default {
  components: {
    modal: Modal,
    'part-form': ParticipantForm,
  },
  emits: ['reload'],
  data: function () {
    return {
      editForm: {
        id: "",
        first_name: "",
        last_name: "",
        bib_number: "",
        birth_year: "",
        team: "",
        gender: "",
      },
    }
  },
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants', 'registerParticipant', 'uploadParticipantCSV', 'updateParticipant']),
    toggle: function () {
      this.$refs.modal.toggle()
    },
    edit: function (participant) {
      this.toggle()

      this.editForm.first_name = participant.firstName
      this.editForm.last_name = participant.lastName
      this.editForm.bib_number = participant.bibNumber
      this.editForm.gender = participant.gender
      this.editForm.team = participant.team
      this.editForm.birth_year = participant.birthYear
      this.editForm.id = participant.id
    },
    reload: function () {
      this.$emit('reload')
    },
    reset: function () {
      this.editForm.first_name = ""
      this.editForm.last_name = ""
      this.editForm.bib_number = ""
      this.editForm.gender = ""
      this.editForm.team = ""
      this.editForm.birth_year = ""
      this.editForm.id = ""
    },
    submit: async function () {
      let success = await this.updateParticipant(this.$route.params.raceId, this.editForm.id, this.editForm)

      if (success) {
        this.toggle()
        this.reload()
      } else {
        console.log("fail")
      }
    },
  }
}
</script>
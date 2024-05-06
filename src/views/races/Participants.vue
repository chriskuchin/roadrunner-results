<template>
  <div class="container mx-5">
    <h1 class="title">Participants</h1>
    <div class="table-container">
      <table class="table is-hoverable is-fullwidth" style="margin: 0 auto;">
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
          <tr v-for="participant in currentParticipants" class="is-size-4">
            <td>{{ participant.bibNumber }}</td>
            <td>{{ participant.firstName }} {{ participant.lastName }}</td>
            <td>{{ participant.birthYear }}</td>
            <td>{{ participant.gender }}</td>
            <td>{{ participant.team }}</td>
            <td><button class="button is-small is-light is-info" @click="openEditModal(participant)">Edit</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <fab @click="openRegisterModal" v-if="isLoggedIn" />
    <edit-participant-modal ref="edit-participant-modal" @reload="reload" />
    <register-participant-modal ref="register-participant-modal" @reload="reload" />
  </div>
</template>

<script>
import FAB from '../../components/Fab.vue'

import { mapActions, mapState } from 'pinia'
import { useUserStore } from '../../store/user'
import { useParticipantsStore } from '../../store/participants'
import EditParticipantModal from '../../components/modals/EditParticipantModal.vue'
import RegisterParticipantModal from '../../components/modals/RegisterParticipantModal.vue'

export default {
  components: {
    'fab': FAB,
    'edit-participant-modal': EditParticipantModal,
    'register-participant-modal': RegisterParticipantModal
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
    }
  },
  mounted: function () {
    this.loadParticipants(this.$route.params.raceId, this.filters.gender, this.filters.team, this.filters.team, this.filters.name, this.filters.limit, this.filters.offset)
  },
  methods: {
    ...mapActions(useParticipantsStore, ['loadParticipants', 'registerParticipant', 'uploadParticipantCSV', 'updateParticipant']),
    ...mapActions(useUserStore, ['isLoggedIn']),
    openRegisterModal: function () {
      this.$refs['register-participant-modal'].toggle()
    },
    reload: function () {
      this.loadParticipants(this.$route.params.raceId, this.filters.gender, this.filters.team, this.filters.team, this.filters.name, this.filters.limit, this.filters.offset)
    },
    openEditModal: function (participant) {
      this.$refs['edit-participant-modal'].edit(participant)
    },
  },
  computed: {
    ...mapState(useParticipantsStore, ['currentParticipants']),
    raceID: function () {
      return this.$route.params.raceId
    },
  }
}
</script>
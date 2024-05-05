<template>
  <div class="section">
    <h1 class="title">Volunteers</h1>
    <tbl :headers="tableHeader" :rows="volunteerInfo" class="mx-auto" />
    <add-volunteer-modal ref="add-volunteer-modal" />
    <fab @click="addVolunteersModal" @close="loadVolunteers" />
  </div>
</template>

<script>
import Fab from '../../components/Fab.vue';
import Table from '../../components/Table.vue';
import AddVolunteersModal from '../../components/modals/AddVolunteers.vue'
import { useRaceStore } from '../../store/race';
import { useErrorBus } from '../../store/error';
import { mapActions, mapStores } from 'pinia';

export default {
  components: {
    'fab': Fab,
    'tbl': Table,
    'add-volunteer-modal': AddVolunteersModal,
  },
  mounted: function () {
    this.loadVolunteers()
  },
  data: function () {
    return {
      tableHeader: [
        {
          abbr: "email",
          title: "Email",
        },
      ],
    }
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    ...mapActions(useRaceStore, ['loadVolunteers']),
    addVolunteersModal: function () {
      this.$refs['add-volunteer-modal'].open(this.$route.params.raceId)
    },
  },
  computed: {
    ...mapStores(useRaceStore),
    volunteerInfo: function () {
      return this.raceStore.volunteers.map((volunteer) => {
        return [
          volunteer.email
        ]
      })
    },
  }
};
</script>
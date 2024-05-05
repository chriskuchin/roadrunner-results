<template>
  <div class="section">
    <h1 class="title">Volunteers</h1>
    <tbl :headers="tableHeader" :rows="volunteerInfo" class="mx-auto" />
    <fab />
  </div>
</template>

<script>
import Fab from '../components/Fab.vue';
import Table from '../components/Table.vue';
import { useRaceStore } from '../store/race';
import { useErrorBus } from '../store/error';
import { mapActions, mapStores } from 'pinia';

export default {
  components: {
    'fab': Fab,
    'tbl': Table,
  },
  mounted: function () {
    this.loadVolunteers()
  },
  data: function () {
    return {
      tableHeader: [
        {
          abbr: "eml",
          title: "Email",
        },
      ],
    }
  },
  methods: {
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    ...mapActions(useRaceStore, ['loadVolunteers']),
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
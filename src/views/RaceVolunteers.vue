<template>
  <div class="section">
    <h1 class="title">Volunteers</h1>
    <ul>
      <li v-for="volunteer in raceStore.volunteers">{{ volunteer.email }}</li>
    </ul>
    <fab />
  </div>
</template>

<script>
import Fab from '../components/Fab.vue';
import { useRaceStore } from '../store/race';
import { mapActions, mapStores } from 'pinia';

export default {
  components: {
    'fab': Fab,
  },
  mounted: function () {
    this.loadVolunteers();
  },
  data: function () {
    return {
    }
  },
  methods: {
    ...mapActions(useRaceStore, ['loadVolunteers']),
  },
  computed: {
    ...mapStores(useRaceStore),
  }
};
</script>
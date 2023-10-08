<template>
  <div class="section">
    <h1 class="title">Volunteers</h1>
    <ul>
      <li v-for="volunteer in volunteers">{{ volunteer.email }}</li>
    </ul>
    <fab />
  </div>
</template>

<script>
import { getRaceVolunteers } from '../api/races'
import Fab from '../components/Fab.vue';

export default {
  components: {
    'fab': Fab,
  },
  data: function () {
    return {
      volunteers: []
    }
  },
  mounted: async function () {
    this.volunteers = await getRaceVolunteers(this.$route.params.raceId)
  }
};
</script>
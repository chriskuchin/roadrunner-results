<template>
  <div class="section">
    <div class="box" v-for="race in racesStore.getRaces" :key="race.id">
      <div class="title is-4">{{ race.name }}</div>
      <router-link :to="'/races/' + race.id + '/events'" class="button">Events</router-link>
    </div>
    <div class="fixed-bottom">
      <a class="button is-primary is-large fab" href="#">
        <icon icon="fa-solid fa-plus"></icon>
      </a>
    </div>
  </div>
</template>

<script>
import { mapStores } from 'pinia'
import { useRacesStore } from "../store/races";

export default {
  components: {
  },
  data: function () {
    return {}
  },
  computed: {
    ...mapStores(useRacesStore),
    raceLink: function (id) {
      console.log(id + " raceLINK")
      return "/races/" + id + "/results"
    }
  },
  mounted: function () {
    this.racesStore.loadRaces()
  }
};
</script>
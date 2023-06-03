<template>
  <div class="container">
    <div class="section mt-1 mb-0 py-0" v-if="this.$route.name != 'events'">
      <router-link :to="'/races/' + this.$route.params.raceId + '/events'">
        <icon icon="fa-solid fa-arrow-left-long"></icon> Back to Events
      </router-link>
    </div>
    <router-view></router-view>
  </div>
</template>

<script>
import { useRaceStore } from "../store/race";

const store = useRaceStore()
export default {
  components: {},
  data: function () {
    return {
    }
  },
  mounted: function () {
    store.loadRace(this.$route.params.raceId)
  },
  methods: {},
  computed: {
    backPath: function () {
      return this.$route.path.slice(0, this.$route.path.lastIndexOf('/'))
    },
  },
};
</script>
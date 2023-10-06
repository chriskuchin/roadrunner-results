<template>
  <div class="container">
    <nav class="breadcrumb mt-3 is-centered mb-0" aria-label="breadcrumbs">
      <ul>
        <li :class="{ 'is-active': !isEventPage }">
          <router-link :to="'/races/' + this.$route.params.raceId + '/events'">
            {{ getName }}
          </router-link>
        </li>
        <li :class="{ 'is-active': isEventPage }" v-if="isEventPage" class="is-active">
          <a href="#" aria-current="page">
            {{ eventName($route.params.eventId) }}
          </a>
        </li>
      </ul>
    </nav> <router-view></router-view>
  </div>
</template>

<script>
import { useRaceStore } from "../store/race";
import { mapState } from "pinia";

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
    ...mapState(useRaceStore, ["getName", "eventName"]),
    backPath: function () {
      return this.$route.path.slice(0, this.$route.path.lastIndexOf('/'))
    },
    isEventPage: function () {
      return this.$route.params.eventId && this.$route.params.eventId != ""
    }
  },
};
</script>
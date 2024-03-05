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
            {{ eventName(getEventId) }}
          </a>
        </li>
      </ul>
    </nav>
    <router-view></router-view>
  </div>
</template>

<script>
import { useRaceStore } from "../store/race";
import { mapState, mapActions } from "pinia";

export default {
  components: {},
  data: function () {
    return {
    }
  },
  mounted: function () {
    this.loadRace(this.$route.params.raceId)
  },
  methods: {
    ...mapActions(useRaceStore, ['loadRace']),
  },
  computed: {
    ...mapState(useRaceStore, ["getName", "eventName"]),
    getEventId: function () {
      if (this.$route.params.eventId && this.$route.params.eventId !== "")
        return this.$route.params.eventId
      else if (this.$route.query.eventId && this.$route.query.eventId !== "")
        return this.$route.query.eventId
    },
    backPath: function () {
      return this.$route.path.slice(0, this.$route.path.lastIndexOf('/'))
    },
    isEventPage: function () {
      return (this.$route.params.eventId && this.$route.params.eventId !== "") || (this.$route.query.eventId && this.$route.query.eventId !== "")
    }
  },
};
</script>
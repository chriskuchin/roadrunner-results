<template>
  <div>
    <div class="box" v-for="event in raceStore.eventList" :key="event.id">
      <div class="title is-4">{{ event.description }}</div>
      <div class="field has-addons">
        <p class="control">
          <router-link :to="getResultsLink(event)" class="button is-primary">Results</router-link>
        </p>
        <p class="control">
          <router-link :to="getTimerLink(event)" class="button">Timer</router-link>
        </p>
        <p class="control">
          <router-link :to="getRecordLink(event)" class="button">Record</router-link>
        </p>
      </div>
    </div>
    <div class="fixed-bottom">
      <a class="button is-primary is-large fab" href="#">
        <icon icon="fa-solid fa-plus"></icon>
      </a>
    </div>
  </div>
</template>

<script>
import { mapStores } from "pinia";
import { useRaceStore } from "../store/race";
export default {
  methods: {
    getResultsLink: function (event) {
      return this.getBaseEventLink(event) + '/results'
    },
    getTimerLink: function (event) {
      return this.getBaseEventLink(event) + "/timer"
    },
    getRecordLink: function (event) {
      return this.getBaseEventLink(event) + "/record"
    },
    getBaseEventLink: function (event) {
      return '/races/' + this.$route.params.raceId + '/events/' + event.eventId
    }
  },
  computed: {
    ...mapStores(useRaceStore),
  }
};
</script>
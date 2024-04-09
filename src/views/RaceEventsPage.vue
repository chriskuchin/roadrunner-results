<template>
  <div class="section">
    <div class="box mb-4" v-for="event in raceStore.eventList " :key="event.id">
      <div class="has-text-right">

        <cm class="is-right" v-if="isLoggedIn">
          <a class="dropdown-item" @click="navigateToAttemptsRecorder(event)" v-if="isDistanceEvent(event)">Attempts</a>
          <a class="dropdown-item" @click="navigateToDiaplayBoard(event)" v-if="isTimerEvent(event)">Display Board</a>
          <a class="dropdown-item" @click="navigateToLaneAssignment(event)" v-if="isTimerEvent(event)">Assign Lanes</a>
          <router-link class="dropdown-item" :to="getTimerLink(event)" v-if="isTimerEvent(event)">Timer</router-link>
          <router-link class="dropdown-item" :to="getRecordLink(event)"
            v-if="isTimerEvent(event)">Recorder</router-link>
          <hr class="dropdown-divider" v-if="isTimerEvent(event)" />
          <a href="#" class="dropdown-item" @click="deleteEvent(event)">Delete Event</a>
        </cm>
      </div>

      <div class="title is-4 mb-4">
        <icon :icon="['fa-solid', eventIcon(event)]"></icon> {{ raceDescription(event) }}
      </div>
      <div class="field has-addons">
        <p class="control">
          <button @click="navigateToEventResults(event)" class="button is-primary">Results</button>
        </p>
        <p class="control">
          <button @click="navigateToDivisionResults(event)" class="button is-link">Division Results</button>
        </p>
      </div>
    </div>
    <fab v-if="isLoggedIn" @click="openModal"></fab>
    <create-event-modal ref="create-event-modal" />
  </div>
</template>

<script>
import { mapStores, mapState, mapActions } from "pinia";
import { useRaceStore } from "../store/race";
import { useUserStore } from "../store/user"
import { useEventStore } from "../store/event";
import { deleteRaceEvent } from "../api/events"
import FAB from '../components/Fab.vue'
import ContextMenu from '../components/DropdownMenu.vue'
import CreateEventModal from "../components/modals/CreateEventModal.vue";

export default {
  components: {
    'create-event-modal': CreateEventModal,
    'fab': FAB,
    'cm': ContextMenu,
  },
  data: function () {
    return {
    }
  },
  methods: {
    ...mapActions(useEventStore, ['loadEvent']),
    navigateToDiaplayBoard(event) {
      this.loadEvent(event)
      return this.$router.push(`${this.getBaseEventLink(event)}/board`)
    },
    navigateToLaneAssignment(event) {
      this.loadEvent(event)
      return this.$router.push(`${this.getBaseEventLink(event)}/lanes`)
    },
    navigateToEventResults(event) {
      this.loadEvent(event)
      return this.$router.push(`${this.getBaseEventLink(event)}/results`)
    },
    navigateToDivisionResults(event) {
      this.loadEvent(event)
      return this.$router.push(`/races/${this.$route.params.raceId}/divisions?eventId=${event.eventId}`)
    },
    navigateToAttemptsRecorder(event) {
      this.loadEvent(event)
      return this.$router.push(`${this.getBaseEventLink(event)}/distance`)
    },
    openModal: function () {
      this.$refs['create-event-modal'].toggle()
    },
    deleteEvent: function (event) {
      var self = this
      let raceID = this.$route.params.raceId
      deleteRaceEvent(raceID, event.eventId).then(() => {
        self.raceStore.loadRace(raceID)
      })
    },
    getResultsLink: function (page, event) {
      return this.getBaseEventLink(event) + `/${page}`
    },
    getTimerLink: function (event) {
      return this.getBaseEventLink(event) + "/timer"
    },
    getRecordLink: function (event) {
      return this.getBaseEventLink(event) + "/record"
    },
    getBaseEventLink: function (event) {
      return '/races/' + this.$route.params.raceId + '/events/' + event.eventId
    },
    isTimerEvent: function (event) {
      return event.type == "relay" || event.type == "timer"
    },
    isDistanceEvent: function (event) {
      return event.type === "distance"
    },
    eventIcon: function (event) {
      switch (event.type) {
        case "timer":
          return 'fa-stopwatch'
        case "distance":
          return 'fa-ruler'
        case "relay":
          return 'fa-repeat'
        default:
          return 'fa-stopwatch'
      }
    },
    raceDescription: function (event) {
      if (event.description == "") {
        return event.distance + "m"
      }

      return event.description

    },
  },
  computed: {
    ...mapState(useUserStore, ['isLoggedIn']),
    raceID: () => this.$route.params.raceId,
    ...mapStores(useRaceStore),
  }
};
</script>../components/DropdownMenu.vue
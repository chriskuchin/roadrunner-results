<template>
  <div class="section">
    <div class="box" v-for="event in  raceStore.eventList " :key="event.id">
      <div class="has-text-right">

        <cm class="is-right" v-if="isLoggedIn">
          <a class="dropdown-item" @click="navigateToAttemptsRecorder(event)" v-if="isDistanceEvent(event)">Attempts</a>
          <a class="dropdown-item" @click="navigateToDiaplayBoard(event)" v-if="isTimerEvent(event)">Display Board</a>
          <router-link class="dropdown-item" :to="getTimerLink(event)" v-if="isTimerEvent(event)">Timer</router-link>
          <router-link class="dropdown-item" :to="getRecordLink(event)"
            v-if="isTimerEvent(event)">Recorder</router-link>
          <hr class="dropdown-divider" v-if="isTimerEvent(event)" />
          <a href="#" class="dropdown-item" @click="deleteEvent(event)">Delete Event</a>
        </cm>
      </div>

      <div class="title is-4">
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
    <fab v-if="isLoggedIn" @click="toggleModal"></fab>
    <modal @close="toggleModal" :show="modal.show">
      <p class="title">Create Event</p>
      <div class="field">
        <label class="label">Description</label>
        <div class="control">
          <input class="input" type="text" placeholder="Event Description" v-model="modal.description">
        </div>
      </div>
      <div class="field">
        <label class="label">Event Type</label>
        <div class="control">
          <div class="select">
            <select v-model="modal.type">
              <option value="timer">Timer</option>
              <option value="distance">Distance</option>
              <option value="relay">Relay</option>
            </select>
          </div>
        </div>
      </div>

      <div class="field has-addons" v-if="modal.type == 'timer'">
        <!-- <label class="label" style="clear:both;">Event Distance</label> -->
        <div class="control">
          <input class="input" type="text" placeholder="Event Distance" v-model="modal.distance_raw">
        </div>
        <div class="control">
          <span class="select">
            <select v-model="modal.distance_unit">
              <option value="meter">Meters</option>
              <option value="kilometer">Kilometer(s)</option>
              <option value="mile">Mile(s)</option>
            </select>
          </span>
        </div>
      </div>

      <div class="field" v-if="modal.type == 'relay'">
        <label class="label">Event Distance</label>
        <div class="control">
          <div class="select">
            <select v-model.number="modal.distance">
              <option value="400">4x100</option>
              <option value="1600">4x400</option>
            </select>
          </div>
        </div>
      </div>

      <div class="field is-grouped">
        <div class="control">
          <button :class="['button', 'is-link', { 'is-loading': modal.creating }]" @click="modalSubmit">Submit</button>
        </div>
        <div class="control">
          <button class="button is-link is-light" @click="toggleModal">Cancel</button>
        </div>
      </div>
    </modal>
  </div>
</template>

<script>
import { mapStores, mapState, mapActions } from "pinia";
import { useRaceStore } from "../store/race";
import { useUserStore } from "../store/user"
import { useEventStore } from "../store/event";
import { createRaceEvent, deleteRaceEvent } from "../api/events"
import Modal from '../components/Modal.vue'
import FAB from '../components/Fab.vue'
import ContextMenu from '../components/DropdownMenu.vue'

export default {
  components: {
    'modal': Modal,
    'fab': FAB,
    'cm': ContextMenu,
  },
  data: function () {
    return {
      modal: {
        show: false,
        creating: false,
        description: "",
        type: "timer",
        distance: 0,
        distance_raw: 1600,
        distance_unit: "meter"
      }
    }
  },
  methods: {
    ...mapActions(useEventStore, ['loadEvent']),
    navigateToDiaplayBoard(event) {
      this.loadEvent(event)
      return this.$router.push(`${this.getBaseEventLink(event)}/board`)
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
    eventDistance: function () {
      if (this.modal.type == "timer") {
        switch (this.modal.distance_unit) {
          case "meter":
            return this.modal.distance_raw
          case "mile":
            return this.modal.distance_raw * 1609.344
          case "kilometer":
            return this.modal.distance_raw * 1000
        }
      } else {
        return this.modal.distance
      }
    },
    eventDescription: function () {
      switch (this.modal.distance_unit) {
        case "meter":
          return `${this.modal.distance_raw}m`
        case "mile":
          if (this.modal.distance_raw == 1)
            return `${this.modal.distance_raw} mile`
          else
            return `${this.modal.distance_raw} miles`
        case "kilometer":
          return `${this.modal.distance_raw}k`
      }

      return `${this.modal.distance_raw} ${this.modal.distance_unit}`
    },
    toggleModal: function () {
      if (!this.modal.show) {
        this.resetModal()
      }

      this.modal.show = !this.modal.show
    },
    modalSubmit: function () {
      let raceID = this.$route.params.raceId
      let eventDistance = this.eventDistance()
      let description = this.modal.description != "" ? this.modal.description : this.eventDescription()

      var self = this
      createRaceEvent(raceID, description, this.modal.type, eventDistance).then(() => {
        self.raceStore.loadRace(raceID)
      })

      this.toggleModal()
    },
    resetModal: function () {
      this.modal.description = ""
      this.modal.type = "timer"
      this.modal.distance = 0
      this.modal.distance_raw = 1600
      this.modal.distance_unit = "meter"
    },
    deleteEvent: function (event) {
      var self = this
      let raceID = this.$route.params.raceId
      deleteRaceEvent(raceID, event.eventId).then(() => {
        self.raceStore.loadRace(raceID)
      })
    },
    modalCancel: function () {
      this.resetModal()
      this.toggleModal()
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
<template>
  <div class="section">
    <div class="box" v-for="event in   raceStore.eventList  " :key="event.id">
      <div class="title is-4">
        <icon :icon="['fa-solid', eventIcon(event)]"></icon> {{ event.description }}
      </div>
      <div class="field has-addons">
        <p class="control">
          <router-link :to="getResultsLink(event)" class="button is-primary">Results</router-link>
        </p>
        <p class="control" v-if="isTimerEvent(event)">
          <router-link :to="getTimerLink(event)" class="button">Timer</router-link>
        </p>
        <p class="control" v-if="isTimerEvent(event)">
          <router-link :to="getRecordLink(event)" class="button">Record</router-link>
        </p>
        <p class="control">
          <a class="button is-danger" @click="deleteEvent(event)">Delete</a>
        </p>
      </div>
    </div>
    <div class="fixed-bottom">
      <a class="button is-primary is-large fab" @click="toggleModal">
        <icon icon="fa-solid fa-plus"></icon>
      </a>
    </div>

    <div :class="['modal', { 'is-active': modal.show }]">
      <div class="modal-background" @click="toggleModal"></div>

      <div class="modal-content">
        <div class="box">
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

          <div class="field" v-if="modal.type == 'timer'">
            <label class="label">Event Distance</label>
            <div class="control">
              <div class="select">
                <select v-model.number="modal.distance">
                  <option value="50">50m</option>
                  <option value="100">100m</option>
                  <option value="200">200m</option>
                  <option value="400">400m</option>
                  <option value="800">800m</option>
                  <option value="1600">1600m</option>
                </select>
              </div>
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
              <button :class="['button', 'is-link', { 'is-loading': modal.creating }]"
                @click="modalSubmit">Submit</button>
            </div>
            <div class="control">
              <button class="button is-link is-light" @click="toggleModal">Cancel</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapStores } from "pinia";
import { useRaceStore } from "../store/race";
import { createRaceEvent, deleteRaceEvent } from "../api/events"

export default {
  data: function () {
    return {
      modal: {
        show: false,
        creating: false,
        description: "",
        type: "timer",
        distance: 1600
      }
    }
  },
  methods: {
    toggleModal: function () {
      // reset the modal before opening it
      if (!this.modal.show) {
        this.resetModal()
      }

      this.modal.show = !this.modal.show
    },
    modalSubmit: function () {
      let raceID = this.$route.params.raceId
      var self = this
      createRaceEvent(raceID, this.modal.description, this.modal.type, this.modal.distance).then(() => {
        self.raceStore.loadRace(raceID)
      })

      this.toggleModal()
    },
    resetModal: function () {
      this.modal.description = ""
      this.modal.type = "timer"
      this.modal.distance = 1600
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
    },
    isTimerEvent: function (event) {
      return event.type == "relay" || event.type == "timer"
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
  },
  computed: {
    raceID: () => this.$route.params.raceId,
    ...mapStores(useRaceStore),
  }
};
</script>
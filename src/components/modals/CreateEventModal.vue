<template>
  <modal ref="create-event-modal" @close="reset">
    <p class="title mb-4">Create Event</p>
    <div class="field mb-2">
      <label class="label">Description</label>
      <div class="control">
        <input class="input" type="text" placeholder="Event Description" v-model="description">
      </div>
    </div>
    <div class="field mb-4">
      <label class="label">Event Type</label>
      <div class="control">
        <div class="select">
          <select v-model="type">
            <option value="timer">Timer</option>
            <option value="distance">Distance</option>
            <option value="relay">Relay</option>
          </select>
        </div>
      </div>
    </div>

    <div class="field has-addons mb-4" v-if="type == 'timer'">
      <!-- <label class="label" style="clear:both;">Event Distance</label> -->
      <div class="control">
        <input class="input" type="text" placeholder="Event Distance" v-model="distance_raw">
      </div>
      <div class="control">
        <span class="select">
          <select v-model="distance_unit">
            <option value="meter">Meters</option>
            <option value="kilometer">Kilometer(s)</option>
            <option value="mile">Mile(s)</option>
          </select>
        </span>
      </div>
    </div>

    <div class="field" v-if="type == 'relay'">
      <label class="label">Event Distance</label>
      <div class="control">
        <div class="select">
          <select v-model.number="distance">
            <option value="400">4x100</option>
            <option value="1600">4x400</option>
          </select>
        </div>
      </div>
    </div>

    <div class="field is-grouped">
      <div class="control">
        <button :class="['button', 'is-link', { 'is-loading': creating }]" @click="submit">Submit</button>
      </div>
      <div class="control">
        <button class="button is-link is-light" @click="toggle">Cancel</button>
      </div>
    </div>
  </modal>
</template>

<script>
import { mapStores } from "pinia";
import { createRaceEvent } from "../../api/events"
import { useRaceStore } from "../../store/race";

import Modal from '../Modal.vue';

export default {
  components: {
    modal: Modal,
  },
  data: function () {
    return {
      show: false,
      creating: false,
      description: "",
      type: "timer",
      distance: 0,
      distance_raw: 1600,
      distance_unit: "meter"
    }
  },
  methods: {
    toggle: function () {
      this.$refs['create-event-modal'].toggle()
    },
    submit: function () {
      let raceID = this.$route.params.raceId
      let eventDistance = this.eventDistance()
      let description = this.description != "" ? this.description : this.eventDescription()

      var self = this
      createRaceEvent(raceID, description, this.type, eventDistance).then(() => {
        self.raceStore.loadRace(raceID)
      })

      this.toggle()
    },
    reset: function () {
      this.description = ""
      this.type = "timer"
      this.distance = 0
      this.distance_raw = 1600
      this.distance_unit = "meter"
    },
    eventDistance: function () {
      if (this.type == "timer") {
        switch (this.distance_unit) {
          case "meter":
            return this.distance_raw
          case "mile":
            return this.distance_raw * 1609.344
          case "kilometer":
            return this.distance_raw * 1000
        }
      } else {
        return this.distance
      }
    },
    eventDescription: function () {
      switch (this.distance_unit) {
        case "meter":
          return `${this.distance_raw}m`
        case "mile":
          if (this.distance_raw == 1)
            return `${this.distance_raw} mile`
          else
            return `${this.distance_raw} miles`
        case "kilometer":
          return `${this.distance_raw}k`
      }

      return `${this.distance_raw} ${this.distance_unit}`
    },
  },
  computed: {
    ...mapStores(useRaceStore),
  }
}
</script>
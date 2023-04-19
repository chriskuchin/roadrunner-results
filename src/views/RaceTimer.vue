<template>
  <div>
    <div>{{ stopwatch }}</div>
    <div class="field has-addons">
      <p class="control">
      <div class="button is-primary" @click="this.startTimer">Start</div>
      </p>
      <p class="control">
      <div class="button is-danger" @click="this.stopTimer">Stop</div>
      </p>
      <p class="control">
      <div class="button is-warning" @click="this.recordFinish">Record</div>
      </p>
    </div>
    {{ timers }}
  </div>
</template>

<script>
import { formatMilliseconds } from "../utilities";

export default {
  components: {
  },
  mounted: function () {
    this.listTimers()
  },
  data: function () {
    return {
      timers: [],
      start: 0,
      duration: 0,
      timer: null
    };
  },
  methods: {
    async listTimers() {
      let res = await fetch("/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/timers")

      this.timers = await res.json()
    },
    async startTimer() {
      let start = Date.now()
      let res = await fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/timers", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          start_ts: start
        })
      })

      if (res.ok) {
        this.start = start
        this.timer = setTimeout(this.tickTimer, 10)
      }
      // POST api/v1/races/:raceId/events/:eventId/timer
      // POST and send timestamp and receive timestamp back
      // feed timestamp into stopwatch display
    },
    stopTimer: function () {
      if (this.timer != null) {
        clearTimeout(this.timer)
        this.timer = null
      }
    },
    async recordFinish() {
      console.log("finisher")

      fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/results", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          end_ts: Date.now()
        })
      })
    },
    tickTimer: function () {
      if (this.start > 0 && this.timer != null) {
        this.duration = Date.now() - this.start
        this.timer = setTimeout(this.tickTimer, 10)
      }
    }
  },
  computed: {
    raceID: function () {
      return this.$route.params.raceId
    },
    eventID: function () {
      return this.$route.params.eventId
    },
    stopwatch: function () {
      return formatMilliseconds(this.duration);
    },
  },
};
</script>
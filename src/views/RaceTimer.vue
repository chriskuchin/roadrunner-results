<template>
  <div>
    Test
    {{ this.$route.params.raceId }}
    {{ this.$route.params.eventId }}
    <div>{{ stopwatch }}</div>
    <div class="button" @click="this.startTimer">Start</div>
    <div class="button" @click="this.stopTimer">Stop</div>
    <div class="button" @click="this.recordFinish">Record</div>
  </div>
</template>

<script>
import { formatStopwatch, formatMilliseconds } from "../utilities";

export default {
  components: {
  },
  data: function () {
    return {
      start: 0,
      duration: 0,
      timer: null
    };
  },
  methods: {
    async startTimer() {
      console.log("StartTimer", this.raceID, this.eventID)
      let res = await fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/timers", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
      })

      if (res.ok) {
        this.start = await res.json()
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
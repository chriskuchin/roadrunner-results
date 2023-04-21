<template>
  <div>
    <div class="level mt-3">
      <div class="level-item has-text-centered">
        <div>
          <p class="title is-1">{{ stopwatch }}</p>
          <p class="heading">Finishers: {{ results }}</p>
        </div>
      </div>
    </div>
    <div class="level mt-1">
      <div class="level-item has-text-centered">
        <div class="field has-addons buttons are-large">
          <p class="control">
          <div class="button is-primary is-responsive" @click="this.startTimer">Start</div>
          </p>
          <p class="control">
          <div class="button is-danger is-responsive" @click="this.stopTimer">Stop</div>
          </p>
          <p class="control">
          <div class="button is-warning is-responsive" @click.passive="this.recordFinish">Record</div>
          </p>
        </div>
      </div>
    </div>

    <div class="tabs is-boxed">
      <ul>
        <li class="is-active"><a>Current Heat</a></li>
        <li v-for="(timer, index) in timers" :key="timer.id" @click="this.clickTab"><a>Heat {{ index + 1 }} ({{
          timer.count }})</a></li>
      </ul>
    </div>
    <div>
      <div v-for="(finisher, index) in finishers" :key="index">
        {{ index + 1 }}. {{ finisher }}
      </div>
    </div>
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
      results: 0,
      timers: [],
      finishers: [],
      start: 0,
      duration: 0,
      timer: null
    };
  },
  methods: {
    clickTab() {
      alert("test")
    },
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
      let finishTime = Date.now()
      fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/results", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          end_ts: finishTime
        })
      })
      this.results++
      this.finishers.push(formatMilliseconds(finishTime - this.start))

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
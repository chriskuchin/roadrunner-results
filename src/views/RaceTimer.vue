<template>
  <div>
    <router-link :to="'/races/' + raceID + '/events'">Back to Events</router-link>
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

    <button class="button is-small is-pulled-right" @click="generateFile">
      <icon icon="fa-solid fa-download"></icon>
    </button>
    <div class="tabs is-boxed">
      <ul>
        <li class="is-active"><a>Current Heat</a></li>
        <li v-for="(timer, index) in timers" :key="timer.id" @click="this.clickTab"><a>Heat {{ index + 1 }} ({{
          timer.count }})</a></li>
      </ul>
    </div>
    <div class="content">
      <ol>
        <li v-for="(finisher, index) in reverseOrderedFinishers" :class="{ unselectable: timerIsRunning }">
          {{ finisher }}
        </li>
      </ol>
    </div>
    <fab @click="fabAction">
      <icon v-if="start == 0" icon="fa-solid fa-play"></icon>
      <icon v-else icon="fa-solid fa-stopwatch"></icon>
    </fab>
    <not :show="error.show" type="is-danger is-light" @close="dismissError">{{ error.msg }}</not>
  </div>
</template>

<script>
import { formatMilliseconds } from "../utilities";
import FAB from '../components/Fab.vue'
import Notification from '../components/Notification.vue'

export default {
  components: {
    'fab': FAB,
    'not': Notification,
  },
  mounted: function () {
    this.listTimers()
  },
  data: function () {
    return {
      reversedFinishers: [],
      results: 0,
      timers: [],
      finishers: [],
      start: 0,
      duration: 0,
      timer: null,
      error: {
        show: false,
        msg: ""
      }
    };
  },
  methods: {
    generateFile() {
      const csvContent = this.finishers.map(row => `${row}\n`).join('');
      const csvData = new Blob([csvContent], { type: 'text/csv' });
      const csvUrl = URL.createObjectURL(csvData);

      const link = document.createElement('a');
      link.href = csvUrl;
      link.download = 'results.csv';

      link.click();
      URL.revokeObjectURL(csvUrl);
    },
    dismissError() {
      this.error.show = false
    },
    fabAction() {
      if (this.start == 0) {
        this.startTimer()
      } else {
        this.recordFinish()
      }
    },
    clickTab() {
      console.log("test")
    },
    async listTimers() {
      let res = await fetch("/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/timers")

      this.timers = await res.json()
    },
    async startTimer() {
      this.start = Date.now()
      this.timer = setTimeout(this.tickTimer, 10)

      let res = await fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/timers", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          start_ts: this.start
        })
      })

      if (res.ok) {
        window.addEventListener("click", this.recordFinish)
      } else {
        this.handleError(`Failed to start the timer: ${res.status}`)
      }
    },
    stopTimer: function () {
      if (this.timer != null) {
        window.removeEventListener("click", this.recordFinish)
        clearTimeout(this.timer)
        this.timer = null
        this.start = 0
        this.duration = 0
      }
    },
    handleError(msg) {
      this.error.msg = msg
      this.error.show = true
    },
    async recordFinish(e) {
      if (e)
        e.stopPropagation()

      let finishTime = Date.now()
      this.results++
      this.finishers.push(formatMilliseconds(finishTime - this.start))

      let res = await fetch(
        "/api/v1/races/" + this.raceID + "/events/" + this.eventID + "/results", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify({
          end_ts: finishTime
        })
      })

      if (!res.ok) {
        this.handleError("Failed to record finisher: " + res.status)
      }
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
    timerIsRunning: function () {
      return this.timer != null
    },
    reverseOrderedFinishers: function () {
      let reversed = []
      for (let i = this.finishers.length; i > 0; i--) {
        reversed.push(this.finishers[i - 1])
      }

      return reversed
    }
  },
};
</script>
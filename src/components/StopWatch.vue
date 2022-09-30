<template>
  <div v-on:keydown="record">
    <h1>TEST</h1>
    {{ stopwatch }}
    <button class="button is-primary" v-on:click="start">Start</button>
    <button class="button" v-on:click="stop">Stop</button>
    <button class="button is-danger" v-on:click="reset">Reset</button>
    <button class="button" v-on:click="record">Record</button>
    {{ results }}
    <button class="button is-warning" v-on:click="clear">Clear</button>
  </div>
</template>

<script>
import { formatStopwatch } from "../utilities";
export default {
  components: {},
  props: [],
  data: function () {
    return {
      timeoutID: null,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
      results: [],
    };
  },
  methods: {
    start: function () {
      this.timeoutID = setTimeout(this.increment, 10);
    },
    stop: function () {
      clearTimeout(this.timeoutID);
    },
    reset: function () {
      this.stop();
      this.timeoutID = null;
      this.milliseconds = 0;
      this.seconds = 0;
      this.minutes = 0;
    },
    record: function () {
      this.results.push({
        timestamp: Date.now(),
        minutes: this.minutes,
        seconds: this.seconds,
        milliseconds: this.milliseconds,
      });
    },
    clear: function () {
      this.results = [];
    },
    increment: function () {
      this.timeoutID = setTimeout(this.increment, 10);

      this.milliseconds++;
      if (this.milliseconds == 100) {
        this.milliseconds = 0;
        this.seconds++;
      }

      if (this.seconds == 60) {
        this.seconds = 0;
        this.minutes++;
      }
    },
  },
  computed: {
    stopwatch: function () {
      return formatStopwatch(this.minutes, this.seconds, this.milliseconds);
    },
  },
};
</script>
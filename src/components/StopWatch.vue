<template>
  <div v-on:keypress="record">
    {{ stopwatch }}
    <div class="field has-addons">
      <p class="control">
        <button class="button is-dark" v-on:click="start">
          <span v-if="running">Pause</span>
          <span v-else>Play</span>
        </button>
      </p>
      <p class="control">
        <button class="button is-dark" v-on:click="reset">Reset</button>
      </p>
      <p class="control">
        <button class="button is-dark" v-on:click="record">Record</button>
      </p>
    </div>
    {{ results }}
    {{ running }}
    <button class="button is-warning" v-on:click="clear">Clear</button>
  </div>
</template>

<script>
import { formatStopwatch } from "../utilities";

import { library } from "@fortawesome/fontawesome-svg-core";
import { faPlay } from "@fortawesome/free-solid-svg-icons";

library.add(faPlay);
export default {
  components: {},
  props: [],
  data: function () {
    return {
      running: false,
      timeoutID: null,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
      results: [],
    };
  },
  methods: {
    start: function () {
      if (!this.running) {
        this.running = true;
        this.timeoutID = setTimeout(this.increment, 10);
      } else {
        this.stop();
      }
    },
    stop: function () {
      this.running = false;
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
    playIcon: function () {
      return faPlay;
    },
  },
};
</script>
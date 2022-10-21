<template>
  <div class="container">
    {{ stopwatch }}
    <div class="field has-addons">
      <p class="control">
        <button class="button is-dark" v-on:click="start">
          <span v-if="running"
            ><font-awesome-icon icon="fa-solid fa-pause" />
          </span>
          <span v-else>
            <font-awesome-icon icon="fa-solid fa-play" />
          </span>
        </button>
      </p>
      <p class="control">
        <button class="button is-dark" v-on:click="reset">
          <font-awesome-icon icon="fa-solid fa-stop" />
        </button>
      </p>
      <p class="control">
        <button class="button is-dark" v-on:click="record" :disabled="!running">
          <font-awesome-icon icon="fa-solid fa-circle" />
        </button>
      </p>
    </div>
    {{ $route.name }}
    <button class="button is-warning" v-on:click="clear">Clear</button>
  </div>
</template>

<script>
import { formatStopwatch, keyHandler } from "../utilities";
import { useResultsStore } from "../store/results";
import { mapActions, mapState } from "pinia";
import { library } from "@fortawesome/fontawesome-svg-core";

import {
  faPlay,
  faPause,
  faStop,
  faCircle,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

library.add(faPlay, faPause, faStop, faCircle);
export default {
  mounted: function () {
    window.addEventListener("keydown", keyHandler(32, this.record));
  },
  unmounted: function () {
    window.removeEventListener("keydown", keyHandler(32, this.record));
  },
  components: {
    "font-awesome-icon": FontAwesomeIcon,
  },
  props: [],
  data: function () {
    return {
      running: false,
      timeoutID: null,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
    };
  },
  methods: {
    ...mapActions(useResultsStore, ["recordFinishTime"]),
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
    record: function (e) {
      if (this.running) {
        console.log("Test Test Test");
        this.recordFinishTime({
          timestamp: Date.now(),
          minutes: this.minutes,
          seconds: this.seconds,
          milliseconds: this.milliseconds,
        });
      }
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
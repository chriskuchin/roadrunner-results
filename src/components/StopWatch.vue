<template>
  <div>
    <h1>TEST</h1>
    {{ stopwatch }}
    <button v-on:click="start">Start</button>
    <button v-on:click="stop">Stop</button>
    <button v-on:click="reset">Reset</button>
  </div>
</template>

<script>
export default {
  components: {},
  props: [],
  data: function () {
    return {
      timeoutID: null,
      minutes: 0,
      seconds: 0,
      milliseconds: 0,
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
      this.timeoutID = null;
      this.milliseconds = 0;
      this.seconds = 0;
      this.minutes = 0;
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
      var minString = this.minutes;
      if (this.minutes < 10) {
        minString = "0" + this.minutes;
      }

      var secString = this.seconds;
      if (this.seconds < 10) {
        secString = "0" + this.seconds;
      }

      var millisString = this.milliseconds;
      if (this.milliseconds < 10) {
        millisString = "0" + this.milliseconds;
      }

      return minString + ":" + secString + ":" + millisString;
    },
  },
};
</script>
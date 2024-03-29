<template>
  <div class="section">
    <div class="level mt-3">
      <video ref="finish-line-camera" id="finish-line-camera" :class="{ active: enableCamera }"></video>
      <div class="level-item has-text-centered">
        <div>
          <p class="title is-1">{{ stopwatch }}</p>
          <p class="heading">Finishers: {{ finisherCount }}</p>
        </div>
      </div>
    </div>
    <div class="level mt-1 mb-1">
      <div class="level-item has-text-centered">
        <div class="field has-addons buttons are-large">
          <p class="control">
          <div v-if="this.timer.id == null" class="button is-primary is-responsive" @click="this.startTimer">Start</div>
          <div v-else class="button is-primary is-responsive" @click="this.resumeTimer">Resume</div>
          </p>
          <p class="control">
          <div class="button is-danger is-responsive" @click="this.stopTimer">Stop</div>
          </p>
          <!-- <p class="control">
          <div class="button is-warning is-responsive" @click.passive="this.recordFinish">Record</div>
          </p> -->
        </div>
      </div>
    </div>
    <!-- <div class="level-item">
      <div class="field">
        <label class="checkbox">
          <input type="checkbox" v-model="enableCamera" @click="manageCamera">
          Finish Line Photos
        </label>
      </div>
    </div> -->

    <button class="button is-small is-pulled-right" @click="generateFile">
      <icon icon="fa-solid fa-download"></icon>
    </button>
    <div class="tabs is-boxed">
      <ul>
        <li :class="{ 'is-active': this.timer.id == null }" @click="this.clickTab(null)"><a>New Heat</a></li>
        <li :class="{ 'is-active': this.timer.id == timer.id }" v-for="(timer, index) in timers" :key="timer.id"
          @click="this.clickTab(timer)">
          <a>
            <icon class="mr-2" v-if="timer.count > 0 && timer.timer_start != 0" icon="fa-solid fa-flag-checkered">
            </icon>
            Heat {{ index + 1 }} ({{ timer.count }})
          </a>
        </li>
      </ul>
    </div>
    <!-- on load of id existing heat load the times for the heat -->
    <div class="content">
      <ol>
        <li v-for="(finisher, index) in reverseOrderedFinishers" :class="{ unselectable: timerIsRunning }">
          {{ finisher }}
        </li>
      </ol>
    </div>
    <fab @click="fabAction">
      <icon v-if="timerStarted()" icon="fa-solid fa-play"></icon>
      <icon v-else icon="fa-solid fa-stopwatch"></icon>
    </fab>
  </div>
</template>

<script>
import { formatMilliseconds } from '../utilities';
import FAB from '../components/Fab.vue'
import { setAuthHeader } from '../api/auth'
import { useErrorBus } from '../store/error';
import { useMediaStore } from '../store/media';
import { mapActions } from 'pinia';
import { listTimers, startTimer } from '../api/timers';
import { recordFinish } from '../api/results';

export default {
  components: {
    'fab': FAB,
  },
  mounted: function () {
    window.addEventListener('keypress', this.handleKeyboardEvent)

    this.listTimers()
    this.loadMedia()
  },
  unmounted: function () {
    window.removeEventListener('keypress', this.handleKeyboardEvent)
  },
  data: function () {
    return {
      timers: [],
      finishers: [],
      enableCamera: false,
      timer: {
        id: null,
        timeout: null,
        start: 0,
        elapsed: 0,
      },
      error: {
        show: false,
        msg: ""
      }
    };
  },
  methods: {
    ...mapActions(useMediaStore, { loadMedia: 'load', startCamera: 'startCamera', stopCamera: 'stopCamera', takePicture: 'takePicture', recordVideo: 'recordVideo', saveVideo: 'saveVideo' }),
    ...mapActions(useErrorBus, { handleError: 'handle' }),
    manageCamera() {
      var video = this.$refs['finish-line-camera']

      if (!this.enableCamera) {
        this.startCamera(video)
      } else {
        this.stopCamera(video)
      }
    },
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
    fabAction() {
      if (this.timerStarted()) {
        this.startTimer()
      } else {
        this.recordFinish()
      }
    },
    clickTab(timer) {
      if (timer) {
        this.timer.id = timer.id
        this.timer.start = timer.timer_start
        this.timer.elapsed = Date.now() - timer.timer_start
      } else {
        this.timer.id = null
        this.timer.start = 0
        this.timer.elapsed = 0
      }
    },
    timerStarted() {
      return this.timer.start == 0
    },
    async listTimers() {
      this.timers = await listTimers(this.raceID, this.eventID)
    },
    resumeTimer() {
      this.timer.timeout = setTimeout(this.tickTimer, 10)
    },
    handleKeyboardEvent(e) {
      if (this.timer.timeout != null && e.keyCode == 32) {
        this.recordFinish()
      }
    },
    async startTimer() {
      window.navigator.vibrate(50)
      this.timer.start = Date.now()
      this.timer.timeout = setTimeout(this.tickTimer, 10)

      try {
        this.timer.id = await startTimer(this.raceID, this.eventID, this.timer.start)
      } catch (err) {
        this.handleError(err)
      }
    },
    stopTimer: function () {
      window.navigator.vibrate(50)
      if (this.timer.timeout != null) {
        this.saveVideo()
        clearTimeout(this.timer.timeout)
        this.timer.timeout = null
        this.timer.start = 0
        this.timer.elapsed = 0

        this.listTimers()
      }
    },
    async recordFinish(e) {
      if (e)
        e.stopPropagation()

      window.navigator.vibrate(50)
      let finishTime = Date.now()
      var elapsedTime = formatMilliseconds(finishTime - this.timer.start)
      this.finishers.push(elapsedTime)
      this.takePicture(this.raceID, this.eventID, finishTime, elapsedTime)

      try {
        await recordFinish(this.raceID, this.eventID, finishTime)
      } catch (err) {
        this.handleError(err)
      }
    },
    tickTimer: function () {
      if (this.timer.start > 0 && this.timer.timeout != null) {
        this.timer.elapsed = Date.now() - this.timer.start
        this.timer.timeout = setTimeout(this.tickTimer, 10)
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
      return formatMilliseconds(this.timer.elapsed);
    },
    timerIsRunning: function () {
      return this.timer.timeout != null
    },
    finisherCount: function () {
      return this.finishers.length
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
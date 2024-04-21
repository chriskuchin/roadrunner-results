<template>
  <div class="mx-auto timer-container">
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
          <div v-if="this.timer.start == 0" class="button is-primary is-responsive is-large" @click="this.startTimer">
            Start
          </div>
          <div v-else class="button is-primary is-responsive is-large" @click="this.resumeTimer">
            Resume
          </div>
          </p>
          <p class="control">
          <div class="button is-danger is-responsive is-large" @click="this.stopTimer">Stop</div>
          </p>
        </div>
      </div>
    </div>
    <div class="tabs is-boxed mt-4">
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
    <div class="results-block">
      <tbl class="mx-auto" :headers="resultsHeader" :rows="tableFinishersPreview" />
    </div>
    <fab @click="fabAction">
      <icon v-if="timerStarted()" icon="fa-solid fa-play"></icon>
      <icon v-else icon="fa-solid fa-stopwatch"></icon>
    </fab>
  </div>
</template>

<style scoped>
.timer-container {
  display: flex;
  flex-direction: column;
  height: 80vh;
  /* Adjust as needed */
  overflow: hidden;
  /* Prevent page from scrolling */
}

.results-block {
  flex: 1;
  /* Grow to fill remaining space */
  overflow-y: auto;
  /* Enable vertical scrolling */
}
</style>


<script>
import { formatMilliseconds } from '../utilities';
import FAB from '../components/Fab.vue'
import Table from '../components/Table.vue';
import { useErrorBus } from '../store/error';
import { useMediaStore } from '../store/media';
import { mapActions } from 'pinia';
import { listTimers, startExistingTimer, startTimer } from '../api/timers';
import { recordFinish } from '../api/results';

export default {
  components: {
    'fab': FAB,
    'tbl': Table,
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
      },
      resultsHeader: [
        {
          abbr: "Pos",
          title: "Position",
        },
        {
          abbr: "Time",
          title: "Finish Time",
        },
        {
          abbr: "Diff",
          title: "Difference From First"
        }
      ],

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
        if (timer.timer_start > 0)
          this.timer.elapsed = Date.now() - timer.timer_start
        else
          this.timer.elapsed = 0
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
        if (this.timer.id !== null) {
          await startExistingTimer(this.raceID, this.eventID, this.timer.id, this.timer.start)
        } else
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
      var elapsedTime = finishTime - this.timer.start
      this.finishers.push(elapsedTime) // needs to be the raw time not formatted
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
    tableFinishersPreview: function () {
      const sourceArray = this.finishers
      if (sourceArray.length === 0) {
        return [];
      }

      // Create an empty array to hold the selected elements and their index
      const resultArray = [];
      const firstTime = this.finishers[0]

      // Calculate the number of elements to retrieve
      const firstThreeCount = Math.min(3, sourceArray.length);
      const lastTenCount = Math.min(5, sourceArray.length - firstThreeCount);

      // Add the first three elements and their index to resultArray
      for (let i = 0; i < firstThreeCount; i++) {
        // <icon v-if="key == 0" icon="fas fa-crown"></icon>
        resultArray.push([`${i + 1}`, formatMilliseconds(sourceArray[i]), formatMilliseconds(sourceArray[i] - firstTime)]);
      }

      if (sourceArray.length > firstThreeCount + lastTenCount)
        resultArray.push(["", "", ""])

      // Add the last ten elements and their index to resultArray
      const startLastTen = sourceArray.length - lastTenCount;
      for (let i = startLastTen; i < sourceArray.length; i++) {
        resultArray.push([i + 1, formatMilliseconds(sourceArray[i]), formatMilliseconds(sourceArray[i] - firstTime)]);
      }
      return resultArray
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
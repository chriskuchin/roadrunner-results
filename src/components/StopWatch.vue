<template>
  <div class="container">
    <div id="scanner-wrapper" ref="scanner-wrapper" :class="{ 'camera-active': cameraActive }">
      <h1>{{ stopwatch }}</h1>
      <video ref="finish-line-camera" id="finish-line-camera" @loadedmetadata="resizeVideo" />
    </div>

    <canvas ref="finish-line-pics" height="300" width="300" id="finish-line-pics" style="display: none"></canvas>
    <div class="field has-addons">
      <p class="control">
        <button class="button is-dark" v-on:click="start">
          <span v-if="running">
            <font-awesome-icon icon="fa-solid fa-pause" />
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
      <p class="control">
        <button class="button is-dark" v-on:click="startCamera">
          start camera
        </button>
      </p>
    </div>
    <button class="button is-warning" v-on:click="clear">Clear</button>
  </div>
</template>

<script>
import { formatStopwatch, keyHandler } from "../utilities";
import { useResultsStore } from "../store/results";
import { mapActions } from "pinia";
import { library } from "@fortawesome/fontawesome-svg-core";

import {
  faPlay,
  faPause,
  faStop,
  faCircle,
} from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import logo from "../assets/images/logo.png";

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
      cameraActive: false,
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
    startCamera: function (e) {
      navigator.mediaDevices
        .getUserMedia({
          audio: false,
          video: {
            mandatory: {
              minWidth: 1280,
              minHeight: 720,
            },
            facingMode: { ideal: "environment" },
          },
        })
        .then((stream) => {
          var video = this.$refs["finish-line-camera"];
          video.srcObject = stream;
          video.play();
          this.cameraActive = true;
        });
    },
    takePic: function (e) {
      var vid = this.$refs["finish-line-camera"];
      var cnvs = this.$refs["finish-line-pics"];
      var ctx = cnvs.getContext("2d");
      console.log(vid.videoHeight, vid.videoWidth);

      ctx.drawImage(vid, 0, 0, 300, 300);
      ctx.font = "20px sans serif";
      ctx.fillStyle = "white";
      ctx.fillText("Hello world", 50, 90);
      var img = new Image();
      img.src = logo;
      img.onload = () => {
        ctx.drawImage(img, 0, 250, 50, 50);
      };

      cnvs.toBlob((blob) => {
        console.log(blob);
      });
    },
    stopCamera: function () {
      const stream = this.$refs["finish-line-camera"].srcObject;
      const tracks = stream.getTracks();

      for (let i = 0; i < tracks.length; i++) {
        let track = tracks[i];
        track.stop();
      }

      this.$refs["finish-line-camera"].srcObject = null;
    },
    resizeVideo: function (e) {
      this.$refs["scanner-wrapper"].style.height = e.target.videoHeight + "px";
      this.$refs["scanner-wrapper"].style.width = e.target.videoWidth + "px";
      console.log(e.target.videoHeight, e.target.videoWidth);
    },
    record: function (e) {
      if (this.running) {
        this.recordFinishTime({
          timestamp: Date.now(),
          minutes: this.minutes,
          seconds: this.seconds,
          milliseconds: this.milliseconds,
        });
      }
    },
    clear: function () { },
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
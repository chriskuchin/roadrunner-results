import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'

export const useMediaStore = defineStore("media", {
  state: () => ({
    videoStream: null,
    videoElement: null,
    recordVideo: false,
    recorder: null,
    chunks: [],
    videoQuality: [
      {
        width: 1920,
        height: 1080,
      },
      {
        width: 1280,
        height: 720,
      },
      {
        width: 960,
        height: 540,
      },
    ],
    settings: {
      audio: false,
      video: {
        width: 1920,
        height: 1080,
        focusMode: 'continuous',
        frameRate: 60,
        facingMode: {
          ideal: "environment"
        }
      }
    }
  }),
  getters: {
  },
  actions: {
    async startCamera(video) {
      this.videoElement = video
      this.videoStream = await navigator.mediaDevices.getUserMedia(this.settings);
      video.srcObject = this.videoStream;
      video.play();
    },
    async recordVideo() {
      if (this.videoStream) {
        const options = {
          videoBitsPerSecond: 2500000000,
          mimeType: "video/webm;codecs=vp8",
        };
        this.recorder = new MediaRecorder(this.videoStream, options)

        var that = this
        this.recorder.ondataavailable = function (e) {
          that.chunks.push(e.data)
        }

        this.recorder.onstop = function () {
          const blob = new Blob(that.chunks, { type: 'video/webm;codecs=vp8' })
          var download = document.createElement("a")
          download.href = URL.createObjectURL(blob)
          download.download = "raceVideo.webm"

          download.click()
          that.chunks = []
          that.recorder = null
        }
        this.recorder.start()
      }
    },
    async saveVideo() {
      if (this.recorder) {
        this.recorder.stop()
      }
    },
    async stopCamera(video) {
      const tracks = this.videoStream.getTracks();

      for (let i = 0; i < tracks.length; i++) {
        let track = tracks[i];
        track.stop();
      }

      video.srcObject = null;
      this.videoStream = null;
    },
    async takePicture(raceID, eventID, finishTime, elapsedTime) {
      if (this.videoStream) {
        const track = this.videoStream.getVideoTracks()[0];
        const canvas = document.createElement("canvas");
        const context = canvas.getContext('2d');

        canvas.width = track.getSettings().width;
        canvas.height = track.getSettings().height;

        context.drawImage(this.videoElement, 0, 0, canvas.width, canvas.height);
        var that = this
        canvas.toBlob(function (imageBlob) {
          that.uploadFinishPhoto(imageBlob, raceID, eventID, finishTime, elapsedTime)
        }, 'image/png')
      }
    },
    async uploadFinishPhoto(imgBlob, raceID, eventID, finishTime, elapsedTime) {
      var url = `/api/v1/races/${raceID}/events/${eventID}/results`
      var formData = new FormData()
      formData.append('photo-finish', imgBlob, `finisher_${finishTime}_${elapsedTime}.png`)

      var res = await fetch(url, await setAuthHeader({
        method: 'POST',
        body: formData
      }))

      if (res.ok) {
        console.log("Success!!!")
      }
    },
    async load() {
      navigator.mediaDevices.enumerateDevices().then(devices => {
        console.log(devices)
      })
    }
  }
})
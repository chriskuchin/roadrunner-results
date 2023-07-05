import { defineStore } from "pinia";
import { setAuthHeader } from '../../api/auth'

export const useMediaStore = defineStore("media", {
  state: () => ({
    videoStream: null,
    videoElement: null,
    recordVideo: false,
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

  }),
  getters: {
  },
  actions: {
    async startCamera(video) {
      this.videoElement = video
      this.videoStream = await navigator.mediaDevices.getUserMedia({
        audio: false,
        video: {
          width: 1920,
          height: 1080,
          facingMode: {
            ideal: "environment"
          }
        }
      });
      video.srcObject = this.videoStream;
      video.play();
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
      if (this.videoStream != null) {
        const track = this.videoStream.getVideoTracks()[0];
        const canvas = document.createElement("canvas");
        const context = canvas.getContext('2d');

        canvas.width = track.getSettings().width;
        canvas.height = track.getSettings().height;

        context.drawImage(this.videoElement, 0, 0, canvas.width, canvas.height);
        console.log("convert image")
        var that = this
        canvas.toBlob(function (imageBlob) {
          console.log("upload image")
          that.uploadFinishPhoto(imageBlob, raceID, eventID, finishTime, elapsedTime)
        }, 'image/png')

        // const imageDataURL = canvas.toDataURL('image/png');

        // const link = document.createElement('a');
        // link.href = imageDataURL;
        // link.download = "finish_x.png";

        // link.click();
      }
    },
    async uploadFinishPhoto(imgBlob, raceID, eventID, finishTime, elapsedTime) {
      var url = `/api/v1/races/${raceID}/events/${eventID}/results`
      var formData = new FormData()
      formData.append('photo-finish', imgBlob, `finisher_${finishTime}_${elapsedTime}.png`)

      var res = await fetch(url, setAuthHeader({
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
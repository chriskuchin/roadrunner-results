<template>
  <div class="section">
    <stopwatch></stopwatch>
    <video ref="finish-line-camera" id="finish-line-camera" />
    <canvas
      ref="finish-line-pics"
      height="300"
      width="300"
      id="finish-line-pics"
    ></canvas>

    <button class="button" @click="startCamera">Start Camera</button>
    <button class="button" @click="stopCamera">Stop Camera</button>
    <button class="button" @click="takePic">Take Picture</button>
  </div>
</template>

<script>
import StopWatch from "../components/StopWatch.vue";
import logo from "../assets/images/logo.png";

export default {
  components: {
    stopwatch: StopWatch,
  },
  data: function () {
    return {
      video: null,
    };
  },
  methods: {
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
    startCamera: function (e) {
      navigator.mediaDevices
        .getUserMedia({
          audio: false,
          video: { facingMode: { ideal: "environment" } },
        })
        .then((stream) => {
          console.log(stream);
          var video = this.$refs["finish-line-camera"];
          video.srcObject = stream;
          video.play();
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
  },
};
</script>
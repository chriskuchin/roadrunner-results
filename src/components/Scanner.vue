<template>
  <div id="qr-code-full-region" style="height: 300px; width: 300px"></div>
</template>

<script>
import { Html5QrcodeScanner } from "html5-qrcode";

export default {
  props: {
    qrbox: {
      type: Number,
      default: 100,
    },
    fps: {
      type: Number,
      default: 10,
    },
  },
  data: function () {},
  mounted: function () {
    const config = {
      fps: this.fps,
      qrbox: this.qrbox,
    };
    const html5QrcodeScanner = new Html5QrcodeScanner(
      "qr-code-full-region",
      config
    );
    html5QrcodeScanner.render(this.onScanSuccess);
  },
  methods: {
    onScanSuccess: function (decodedText, decodedResult) {
      console.log(decodedText, decodedResult);
    },
  },
};
</script>
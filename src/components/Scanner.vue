<template>
  <div id="scanner"></div>
  <button @click="toggleScanner" class="button">Start Scanning</button>
</template>

<script>
import { Html5Qrcode, Html5QrcodeScannerState } from "html5-qrcode";

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
  data: function () {
    return {
      html5QRCode: null,
      qrCodeConfig: {
        fps: this.fps,
        qrbox: { height: this.qrbox, width: this.qrbox },
      },
    };
  },
  mounted: function () {
    this.html5QRCode = new Html5Qrcode("scanner");
  },
  unmounted: function () {
    if (this.html5QRCode.getState() == Html5QrcodeScannerState.SCANNING) {
      this.html5QRCode
        .stop()
        .then((ignore) => {})
        .catch((err) => {
          console.log("Failed to stop", err);
        });
    }
  },
  methods: {
    onScanSuccess: function (decodedText, decodedResult) {
      console.log(decodedText, decodedResult);
    },
    toggleScanner: function () {
      this.html5QRCode.start(
        { facingMode: "environment" },
        this.qrCodeConfig,
        this.onScanSuccess
      );
    },
  },
};
</script>
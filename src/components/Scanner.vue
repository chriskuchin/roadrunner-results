<template>
  <div id="scanner"></div>
  <div class="select">
    <select>
      <option>Select dropdown</option>
      <option>With options</option>
    </select>
  </div>
  <button @click="toggleScanner" class="button">{{ scannerToggleDescription }}</button>
  active: {{ scannerActive }}
  inactive: {{ scannerInactive }}
  unknown: {{ scannerStateUnknown }}
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
    if (this.scannerActive) {
      this.html5QRCode
        .stop()
        .then((ignore) => { })
        .catch((err) => { });
    }
  },
  methods: {
    onScanSuccess: function (decodedText, decodedResult) {
      console.log(decodedText, decodedResult);
    },
    toggleScanner: function () {
      // console.log(this.html5QRCode.getState() == Html5QrcodeScannerState.PAUSED, "paused")
      // console.log(this.html5QRCode.getState() == Html5QrcodeScannerState.SCANNING, "scanning")
      // console.log(this.html5QRCode.getState() == Html5QrcodeScannerState.UNKNOWN, "unknown")
      // console.log(this.html5QRCode.getState() == Html5QrcodeScannerState.NOT_STARTED, "not started")
      if (this.scannerInactive) {
        this.html5QRCode.start(
          { facingMode: "environment" },
          this.qrCodeConfig,
          this.onScanSuccess
        );
      } else {
        this.html5QRCode.stop().then((ignore) => {
          // QR Code scanning is stopped.
        }).catch((err) => {
          // Stop failed, handle it.
        });
      }
    },
  },
  computed: {
    scannerToggleDescription: function () {
      return !this.scannerActive ? "Start Scanner" : "Stop Scanner"
    },
    scannerStateUnknown: function () {
      return this.html5QRCode && this.html5QRCode.getState() == Html5QrcodeScannerState.UNKNOWN
    },
    scannerActive: function () {
      return this.html5QRCode && this.html5QRCode.getState() == Html5QrcodeScannerState.SCANNING
    },
    scannerInactive: function () {
      return this.html5QRCode == null || this.html5QRCode.getState() == Html5QrcodeScannerState.NOT_STARTED
    }
  }
};
</script>
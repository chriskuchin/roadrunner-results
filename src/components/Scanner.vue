<template>
  <div>
    <div id="scanner"></div>
    <button @click="toggleScanner" class="button">
      {{ scannerToggleDescription }}
    </button>
  </div>
</template>

<script>
import { Html5Qrcode, Html5QrcodeScannerState } from "html5-qrcode";

export default {
  props: {
    qrbox: {
      type: Number,
      default: 250,
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
        qrbox: this.qrbox,
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
      let val = decodedText.split("-")
      if (val.length == 2) {
        this.$emit("bib", {
          bib: val[0]
        })
      } else {
        this.$emit("bib", {
          bib: val
        })
      }
    },
    toggleScanner: function () {
      if (this.scannerInactive) {
        this.html5QRCode.start(
          { facingMode: "environment" },
          this.qrCodeConfig,
          this.onScanSuccess
        );
      } else {
        this.html5QRCode
          .stop()
          .then((ignore) => {
            // QR Code scanning is stopped.
          })
          .catch((err) => {
            // Stop failed, handle it.
          });
      }
    },
  },
  computed: {
    scannerToggleDescription: function () {
      return !this.scannerActive ? "Start Scanner" : "Stop Scanner";
    },
    scannerStateUnknown: function () {
      return (
        this.html5QRCode &&
        this.html5QRCode.getState() == Html5QrcodeScannerState.UNKNOWN
      );
    },
    scannerActive: function () {
      return (
        this.html5QRCode &&
        this.html5QRCode.getState() == Html5QrcodeScannerState.SCANNING
      );
    },
    scannerInactive: function () {
      return (
        this.html5QRCode == null ||
        this.html5QRCode.getState() == Html5QrcodeScannerState.NOT_STARTED
      );
    },
  },
};
</script>
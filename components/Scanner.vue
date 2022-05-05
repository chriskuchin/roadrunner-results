<template>
  <div>
    <div id="scanner"></div>
    {{ scannedResults }}
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import { Html5QrcodeScanner, Html5QrcodeScanType } from "html5-qrcode";

export default Vue.extend({
  name: "CodeScanner",
  data() {
    return {
      scannedResults: [],
    };
  },
  mounted() {
    const scanner = new Html5QrcodeScanner(
      "scanner",
      {
        fps: 100,
        qrbox: { width: 250, height: 250 },
        supportedScanTypes: [Html5QrcodeScanType.SCAN_TYPE_CAMERA],
      },
      false
    );
    scanner.render(this.onScan, function (err) {
      console.log(err);
    });
  },
  methods: {
    onScan: function (decodedText, decodedResult) {
      if (this.scannedResults.includes(decodedText)) {
        return;
      }
      this.scannedResults.push(decodedText);
    },
  },
});
</script>
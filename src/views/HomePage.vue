<template>
  <div class="section">
    <div id="scanner"></div>
  </div>
</template>

<script>
import { Html5QrcodeScanner } from "html5-qrcode"

export default {
  components: {},
  mounted: function () {
    var scanner = new Html5QrcodeScanner("scanner", {})
    scanner.render()
  }
};
</script>
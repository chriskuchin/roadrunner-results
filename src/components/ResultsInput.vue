<template>
  <div>
    <input
      class="input is-large"
      type="number"
      placeholder="Bib Number Input"
      v-on:keyup.enter="recordRunner"
    />
  </div>
</template>

<script>
import { mapActions } from "pinia";
import { useResultsStore } from "../store/results";

export default {
  data: function () {
    return {};
  },
  methods: {
    ...mapActions(useResultsStore, ["recordRunnerResult"]),
    recordRunner: function (e) {
      e.preventDefault();
      var bibNumber = e.currentTarget.value;
      e.currentTarget.value = "";

      this.recordRunnerResult({
        bib: bibNumber,
        timestamp: Date.now(),
      });
    },
  },
};
</script>
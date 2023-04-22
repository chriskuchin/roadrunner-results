<template>
  <div>
    <input class="input is-large" :type="getInputType" placeholder="Bib Number Input" v-on:keyup.enter="recordRunner" />
    <label class="checkbox">
      <input type="checkbox" v-model="letterToggle">
      Allow Letters
    </label>
  </div>
</template>

<script>
import { mapActions } from "pinia";
import { useResultsStore } from "../store/results";

export default {
  props: ['raceId', 'eventId', 'timerId'],
  data: function () {
    return {
      letterToggle: false,
    };
  },
  methods: {
    ...mapActions(useResultsStore, ["recordRunnerResult"]),
    recordRunner: function (e) {
      e.preventDefault();
      var bibNumber = e.currentTarget.value;
      e.currentTarget.value = "";

      this.recordRunnerResult({
        bib: bibNumber,
        raceId: this.raceId,
        eventId: this.eventId,
        timerId: this.timerId,
      });
    },
  },
  computed: {
    getInputType() {
      return this.letterToggle ? "text" : "number"
    }
  }
};
</script>
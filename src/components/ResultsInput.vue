<template>
  <div>
    <input class="input is-large is-danger" :type="getInputType" placeholder="Bib Number Input"
      v-on:keyup.enter="recordRunner" />
    <p class="help is-danger">Failed to record runner bib</p>
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
  emits: ['recorded-racer'],
  data: function () {
    return {
      letterToggle: false,
    };
  },
  methods: {
    ...mapActions(useResultsStore, ["recordRunnerResult"]),
    inputError: function() {
      // set error state and error message
      // double buzz?
      // set time out function to clear message
    },
    inputSuccess: function() {
      // set success state and success message maybe the time it's associated with
      // single buzz
      // set timeout function to clear
    },
    recordRunner: async function (e) {
      e.preventDefault();
      var bibNumber = e.currentTarget.value;
      e.currentTarget.value = "";

      await this.recordRunnerResult({
        bib: bibNumber,
        raceId: this.raceId,
        eventId: this.eventId,
        timerId: this.timerId,
      });

      this.$emit('recorded-racer')
    },
  },
  computed: {
    getInputType() {
      return this.letterToggle ? "text" : "number"
    }
  }
};
</script>
<template>
  <div>
    <input :class="['input', 'is-large', result.class]" ref="input" :type="getInputType" placeholder="Bib Number Input"
      v-on:keyup.enter="recordRunner" />
    <p :class="['help', result.class]" v-if="result.show">{{ result.msg }}</p>
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
      result: {
        show: false,
        class: "",
        msg: "",
      }
    };
  },
  methods: {
    ...mapActions(useResultsStore, ["recordRunnerResult"]),
    clearStatus: function () {
      this.result.show = false
      this.result.class = ""
    },
    inputError: function () {
      // set error state and error message
      this.result.class = "is-danger"
      this.result.msg = "Failed to record bib number"
      this.result.show = true

      // double buzz?
      window.navigator.vibrate([50, 10, 50])
    },
    inputSuccess: function () {
      // set success state and success message maybe the time it's associated with
      this.result.class = "is-success"
      this.result.msg = "Success"
      this.result.show = true
      // single buzz
      window.navigator.vibrate([50])

      // set timeout function to clear
      setTimeout(this.clearStatus, 1000);
    },
    recordRunner: async function (e) {
      e.preventDefault();
      var bibNumber = this.$refs.input.value

      let ok = await this.recordRunnerResult({
        bib: bibNumber,
        raceId: this.raceId,
        eventId: this.eventId,
        timerId: this.timerId,
      });

      if (ok) {
        this.$refs.input.value = ""
        this.$emit('recorded-racer')
        console.log("success")
        this.inputSuccess()
      } else {
        console.log("error")
        this.inputError()
      }
    },
  },
  computed: {
    getInputType() {
      return this.letterToggle ? "text" : "number"
    }
  }
};
</script>
<template>
  <div>
    <div class="field has-addons">
      <p class="control">
        <a :class="['button', 'is-static', 'is-large', result.class]">
          {{ finisher }}/{{ totalResults }}
        </a>
      </p>
      <p class="control is-expanded">
        <input :class="['input', 'is-large', result.class]" ref="input" :type="getInputType"
          placeholder="Bib Number Input" v-on:keyup.enter="recordRunner" />
      <p :class="['help', result.class]" v-if="result.show">{{ result.msg }}</p>
      </p>
      <p class="control" v-if="false">
        <a :class="['button', 'is-static', 'is-large', result.class]">
          {{ time }}
        </a>
      </p>
    </div>
    <label class="checkbox">
      <input type="checkbox" v-model="letterToggle">
      Allow Letters
    </label>
  </div>
</template>

<script>
export default {
  props: ['raceId', 'eventId', 'timerId', 'totalResults', 'finisher', 'time'],
  emits: ['bib'],
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
      this.$refs.input.value = ""

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

      this.$emit("bib", {
        bib: bibNumber,
        success: this.inputSuccess,
        error: this.inputError,
      })
    },
  },
  computed: {
    getInputType() {
      return this.letterToggle ? "text" : "number"
    }
  }
};
</script>
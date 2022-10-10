<template>
  <div>
    <result-scanner />
    <result-input />
  </div>
</template>

  <script>
import Scanner from "../components/Scanner.vue";
import RacerInput from "../components/ResultsInput.vue";

export default {
  components: {
    "result-scanner": Scanner,
    "result-input": RacerInput,
  },
  data: function () {
    return {};
  },
  created: function () {},
  computed: {},
};
</script>
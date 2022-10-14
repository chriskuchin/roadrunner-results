<template>
  <div class="section">
    <div class="tabs">
      <ul>
        <li class="is-active"><a>Scanner</a></li>
        <li><a>Manual</a></li>
      </ul>
    </div>
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
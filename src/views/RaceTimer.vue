<template>
  <div class="section">
    <stopwatch></stopwatch>
  </div>
</template>

<script>
import StopWatch from "../components/StopWatch.vue";

export default {
  components: {
    stopwatch: StopWatch,
  },
  data: function () {
    return {
      video: null,
    };
  },
  methods: {
  },
};
</script>
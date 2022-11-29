<template>
  <div class="section">
    <breadcrumb></breadcrumb>
    <stopwatch></stopwatch>
  </div>
</template>

<script>
import StopWatch from "../components/StopWatch.vue";
import Breadcrumb from "../components/Breadcrumb.vue";
export default {
  components: {
    stopwatch: StopWatch,
    breadcrumb: Breadcrumb,
  },
  data: function () {
    return {
      video: null,
    };
  },
  methods: {},
};
</script>
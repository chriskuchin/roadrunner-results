<template>
  <div>
    <header>
      <breadcrumb :paths="paths"></breadcrumb>
    </header>
    <div class="section">
      <stopwatch></stopwatch>
    </div>
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
  computed: {
    paths: function () {
      return [
        {
          path: "/races",
          display: "Races",
        },
        {
          path: "/races/" + this.$route.params.raceId,
          display: "Race 123",
        },
        {
          path: "/races/" + this.$route.params.raceId + "/timer",
          display: "Timer",
        },
      ];
    },
  },
};
</script>
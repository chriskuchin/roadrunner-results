<template>
  <div>
    <header>
      <breadcrumb :paths="paths"></breadcrumb>
    </header>

    Hello World Register a Racer
  </div>
</template>

<script>
import Breadcrumb from "../components/Breadcrumb.vue";

export default {
  components: {
    breadcrumb: Breadcrumb,
  },
  data: function () {
    return {};
  },
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
          path: "/races/" + this.$route.params.raceId + "/register",
          display: "Register",
        },
      ];
    },
  },
};
</script>
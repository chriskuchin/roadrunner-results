<template>
  <breadcrumb :paths="paths"></breadcrumb>

  Hello World Register a Racer
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
  computed: function () {
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
};
</script>
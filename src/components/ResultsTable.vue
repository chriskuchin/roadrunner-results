<template>
  <div>
    {{ mergedResults }}
  </div>
</template>

<script>
import { mapState } from "pinia";
import { useResultsStore } from "../store/results";

export default {
  props: {},
  data: function () {
    return {};
  },
  computed: {
    ...mapState(useResultsStore, ["mergedResults"]),
  },
};
</script>
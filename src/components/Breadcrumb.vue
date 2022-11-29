<template></template>

<script>
export default {
  mounted: function () {
    console.log(this.$route);
    console.log(this.$route.name);
    console.log(this.$route.params.raceId);
  },
};
</script>
<template>
  <nav class="breadcrumb" aria-label="breadcrumb">
    <ul>
      <li v-for="crumb in paths" :key="crumb.path">
        <router-link :to="crumb.path">{{ crumb.display }}</router-link>
      </li>
    </ul>
  </nav>
</template>

<script>
export default {
  props: ["paths"],
  data: function () {
    return {};
  },
  mounted: function () {
    console.log(this.$route);
    console.log(this.$route.name);
    console.log(this.$route.params.raceId);
  },
};
</script>
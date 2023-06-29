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
  props: ["paths"]
};
</script>
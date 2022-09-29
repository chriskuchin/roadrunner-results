<template>
  <div id="app">
    <site-navigation></site-navigation>
    <router-view></router-view>
  </div>
</template>

<script>
import SiteNavigation from './components/SiteNavigation.vue'

export default {
    components: {
        "site-navigation": SiteNavigation
    }
}

</script>
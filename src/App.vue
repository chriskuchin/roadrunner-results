<template>
  <div id="app">
    <site-navigation></site-navigation>
    <router-view></router-view>
    <notification :show="showNotification" type="is-danger is-light" @close="dismissError">{{ errorMsg }}</notification>
  </div>
</template>

<script>
import SiteNavigation from './components/SiteNavigation.vue'
import { mapActions, mapState } from 'pinia'
import { useErrorBus } from './store/error'
import Notification from './components/Notification.vue'

export default {
  components: {
    "site-navigation": SiteNavigation,
    "notification": Notification
  },
  computed: {
    ...mapState(useErrorBus, ['showNotification', 'errorMsg'])
  },
  methods: {
    ...mapActions(useErrorBus, { dismissError: 'hide' })
  }
}

</script>
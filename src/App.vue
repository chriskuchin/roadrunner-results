<template>
  <div id="app">
    <site-navigation></site-navigation>
    <router-view></router-view>
    <notification :show="hasError" type="is-danger is-light" @close="dismissError">
      {{ errorMsg }} (1/{{ errorCount }})
    </notification>
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
    ...mapState(useErrorBus, ['hasError', 'errorMsg', 'errorCount'])
  },
  methods: {
    ...mapActions(useErrorBus, { dismissError: 'hide' })
  }
}

</script>
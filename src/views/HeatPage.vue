<template>
  <div class="container">
    <fab></fab>
  </div>
</template>

<script>
import { setAuthHeader } from '../api/auth';
import FAB from '../components/Fab.vue'

export default {
  components: {
    "fab": FAB,
  },
  mounted: async function () {
  },
  data: function () {
    return {
    }
  },
  computed: {
  },
  methods: {
  }
};
</script>
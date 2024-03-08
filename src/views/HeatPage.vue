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
    await fetch(`/api/v1/races/59a8fa78-95e5-4cb2-befa-6264e41bb797/events/87d00d7d-e6cb-4570-bcf7-1c48a9aadf72/heats`, await setAuthHeader({
      method: "POST",
      body: JSON.stringify({
        "assignments": [
          {
            "lane": 1,
            "bib": "331"
          },
          {
            "lane": 3,
            "bib": "231"
          }
        ]
      }),
    }))
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
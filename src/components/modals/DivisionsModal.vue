<template>
  <modal ref="divisions-modal">
    <div v-for="division in divisions">{{ division.display }} - {{ division.filters }}</div>
  </modal>
</template>

<script>
import { mapStores, mapState, mapActions } from 'pinia'
import { useDivisionsStore } from '../../store/divisions'

import Modal from '../Modal.vue';

export default {
  components: {
    modal: Modal,
  },
  data: function () {
    return {
      raceId: ""
    }
  },
  methods: {
    ...mapActions(useDivisionsStore, ['createDivision', 'load']),
    toggle: function () {
      this.$refs['divisions-modal'].toggle()
    },
    open: function (raceId) {
      this.toggle()
      this.raceId = raceId
      this.load(this.raceId).then(() => {
        console.log("Loaded Divisions")
      })
    },
  },
  computed: {
    ...mapState(useDivisionsStore, ['divisions']),
  }

}
</script>
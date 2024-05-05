<template>
  <modal ref="add-volunteer-modal" @close="reset">
  </modal>
</template>

<script>
import { mapActions } from "pinia";
import { useRaceStore } from "../../race";
import Modal from '../Modal.vue';

export default {
  components: {
    modal: Modal,
  },
  methods: {
    ...mapActions(useRaceStore, ['loadEvents']),
    toggle: function () {
      this.$refs['add-volunteer-modal'].toggle()
    },
  }
}
</script>
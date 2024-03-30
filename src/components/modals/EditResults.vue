<template>
  <modal ref="modal">
    <div class="field is-horizontal">
      <div class="field-label is-normal">
        <label class="label">From</label>
      </div>
      <div class="field-body">
        <div class="field">
          <p class="control">
            <input class="input" type="number" placeholder="Time">:
          </p>
        </div>
        <div class="field">
          <p class="control is-expanded has-icons-left has-icons-right">
            <input class="input is-success" type="email" placeholder="Bib Number" value="">
            <span class="icon is-small is-left">
              <i class="fas fa-envelope"></i>
            </span>
            <span class="icon is-small is-right">
              <i class="fas fa-check"></i>
            </span>
          </p>
        </div>
      </div>
    </div>
  </modal>
</template>

<script>
import Modal from '../Modal.vue'

export default {
  components: {
    modal: Modal,
  },
  data: function () {
    return {}
  },
  methods: {
    open: function () {
      this.$refs.modal.toggle()
    },
    close: function () {
      this.$refs.modal.toggle()
    }
  },
}
</script>
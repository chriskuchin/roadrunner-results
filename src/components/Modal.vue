<template>
  <div ref="modal" class="modal">
    <div class="modal-background" @click="toggle"></div>
    <div class="modal-content">
      <div class="box">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  emits: ["close", "open"],
  methods: {
    toggle(e) {
      const el = this.$refs.modal

      el.classList.toggle('is-active')
      if (el.classList.contains("is-active")) {
        this.$emit('open', e)
      } else {
        this.$emit('close', e)
      }
    },
  }
};
</script>
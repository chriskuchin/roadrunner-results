<template>
  <div :class="['modal', { 'is-active': show }]">
    <div class="modal-background" @click="$emit('close')"></div>
    <div class="modal-content">
      <div class="box">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  emits: ["close"],
  props: {
    show: Boolean
  },
};
</script>
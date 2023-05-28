<template>
  <div :class="['notification', type, 'is-bottom-slide', { 'show': show }]">
    <button @click="$emit('close')" class="delete"></button>
    <slot>test</slot>
  </div>
</template>


<script>
export default {
  emits: ["close"],
  props: {
    "show": Boolean,
    "type": String,
  }
};
</script>

<style>
.notification.is-bottom-slide {
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  transform: translateY(100%);
  transition: transform 0.3s ease-out;
}

.notification.is-bottom-slide.show {
  transform: translateY(0);
}
</style>
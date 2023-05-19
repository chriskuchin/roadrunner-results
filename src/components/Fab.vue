<template>
  <div class="fixed-bottom">
    <a class="button is-primary is-large fab" @click="$emit('click')">
      <slot>
        <icon icon="fa-solid fa-plus"></icon>
      </slot>
    </a>
  </div>
</template>

<script>
export default {
  emits: ["click"]
};
</script>
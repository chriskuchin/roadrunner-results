<template>
  <div class="fixed-bottom">
    <a :class="['button', button_type, 'is-large', 'fab']" @click="emitClick">
      <slot>
        <icon icon="fa-solid fa-plus"></icon>
      </slot>
    </a>
  </div>
</template>

<script>
export default {
  props: {
    is_disabled: {
      type: Boolean,
      default: false
    },
    button_type: {
      type: String,
      default: "is-primary"
    }
  },
  emits: ["click"],
  computed: {
    button_disabled: function () { }
  },
  methods: {
    emitClick() {
      this.$emit('click')
    }
  }
};
</script>
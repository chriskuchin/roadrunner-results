<template>
  <div class="dropdown" @click="toggleMenu">
    <div class="dropdown-trigger">
      <span class="icon is-clickable" aria-haspopup="true" aria-controls="dropdown-menu">
        <icon icon="fa-solid fa-ellipsis-v"></icon>
      </span>
    </div>
    <div class="dropdown-menu" role="menu">
      <div class="dropdown-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
  },
  emits: [],
  computed: {
  },
  methods: {
    toggleMenu: function (e) {
      let target = e.currentTarget
      if (target.className.includes("is-active")) {
        target.className = target.className.replace("is-active", "")
      }
      else {
        target.className = e.currentTarget.className + " is-active"
      }
    },
  }
};
</script>